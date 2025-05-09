package repo

import (
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"code.gitea.io/gitea/models/db"
	git_model "code.gitea.io/gitea/models/git"
	issue_model "code.gitea.io/gitea/models/issues"
	access_model "code.gitea.io/gitea/models/perm/access"
	repo_model "code.gitea.io/gitea/models/repo"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/optional"
	"code.gitea.io/gitea/modules/setting"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/modules/util"
	"code.gitea.io/gitea/services/context"
	"code.gitea.io/gitea/services/convert"
	issue_service "code.gitea.io/gitea/services/issue"
	repo_service "code.gitea.io/gitea/services/repository"
	security_service "code.gitea.io/gitea/services/security"
)

const (
	tplSecurityContainerScanning       base.TplName = "repo/security/container_scanning"
	tplSecurityDetailContainerScanning base.TplName = "repo/security/detail_container_scanning"
	tplSecurityStart                   base.TplName = "repo/security/start"
	tplSecurityDependencyList          base.TplName = "repo/security/dependency_list"
	tplSecurityDependencyScanning      base.TplName = "repo/security/dependency_scanning"
	tplSecurityDependencyVuln          base.TplName = "repo/security/dependency_vuln"
	tplSecurityIaCMisconfig            base.TplName = "repo/security/iac_misconfiguration"
	tplSecurityDetailIaCMisconfig      base.TplName = "repo/security/detail_iac_misconfiguration"
	tplSecuritySecretDetection         base.TplName = "repo/security/secret_detection"
	tplSecurityDetailSecretDetection   base.TplName = "repo/security/detail_secret_detection"
	tplDashboard                       base.TplName = "repo/security/security_dashboard"
	tplSecurityCodeScanning            base.TplName = "repo/security/code_scanning"
	tplSecurityCodeVuln                base.TplName = "repo/security/code_vuln"
)

// Security renders single security page
func Security(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.security")
	ctx.Data["PageIsSecurity"] = true
	ctx.Data["PageIsOverview"] = true

	PrepareSecurityBranchOption(ctx)

	ctx.HTML(http.StatusOK, tplSecurityStart)
}

func PrepareSecurityBranchOption(ctx *context.Context) {
	branchOpts := git_model.FindBranchOptions{
		RepoID:          ctx.Repo.Repository.ID,
		IsDeletedBranch: optional.Some(false),
		ListOptions:     db.ListOptionsAll,
	}
	brs, err := git_model.FindBranchNames(ctx, branchOpts)
	if err != nil {
		ctx.ServerError("GetBranches", err)
		return
	}
	// always put default branch on the top if it exists
	if slices.Contains(brs, ctx.Repo.Repository.DefaultBranch) {
		brs = util.SliceRemoveAll(brs, ctx.Repo.Repository.DefaultBranch)
		brs = append([]string{ctx.Repo.Repository.DefaultBranch}, brs...)
	}
	ctx.Data["Branches"] = brs
}

// DependencyList renders single depedency list page
func DependencyList(ctx *context.Context) {
	repo := ctx.Repo.Repository
	ctx.Data["Title"] = ctx.Tr("repo.security")

	page := ctx.FormInt("page")
	if page <= 1 {
		page = 1
	}
	dependencyListFiltered := map[string]string{
		"location": "",
		"licenses": "",
		"q":        "",
	}
	for filterKey := range dependencyListFiltered {
		if filterValue := ctx.FormTrim(filterKey); filterValue != "" {
			dependencyListFiltered[filterKey] = filterValue
		}
	}

	allDependencyList, _ := repo_model.GetDependencyList(ctx, repo.ID, dependencyListFiltered)
	total := len(allDependencyList)

	dependencyList, _, err := db.FindAndCount[repo_model.RepoDependency](ctx, repo_model.RepoDependencySearchOptions{
		RepoDependency: repo_model.RepoDependency{
			RepoID:   repo.ID,
			Target:   dependencyListFiltered["location"],
			Licenses: dependencyListFiltered["licenses"],
		},
		ListOptions: db.ListOptions{
			PageSize: setting.UI.IssuePagingNum,
			Page:     page,
		},
		Q: dependencyListFiltered["q"],
	})

	if err != nil {
		ctx.ServerError("GetDependencyList", err)
		return
	}

	apiDependencyList := make([]*api.RepoDependency, len(dependencyList))
	for i := range dependencyList {
		apiDependencyList[i] = convert.ToRepoDependency(ctx, dependencyList[i])
	}

	pager := context.NewPagination(total, setting.UI.IssuePagingNum, page, 10)
	pager.AddParamString("q", dependencyListFiltered["q"])
	pager.AddParamString("location", dependencyListFiltered["location"])
	pager.AddParamString("licenses", dependencyListFiltered["licenses"])
	ctx.Data["Page"] = pager

	ctx.Data["PageIsSecurity"] = true
	ctx.Data["PageIsDependencyList"] = true
	ctx.Data["DependencyListFilters"] = repo_model.GetListDependencyFilter(ctx, repo.ID)
	ctx.Data["PaginatedDependencyList"] = apiDependencyList
	ctx.Data["DependencyListFiltered"] = dependencyListFiltered
	ctx.HTML(http.StatusOK, tplSecurityDependencyList)
}

// DepedencyScanning renders single depedency vulnerabilities page
func DepedencyScanning(ctx *context.Context) {
	repo := ctx.Repo.Repository
	ctx.Data["Title"] = ctx.Tr("repo.security")

	page := ctx.FormInt("page")
	if page < 1 {
		page = 1
	}
	itemsPerPage := 10

	vulns, err := repo_model.ListDependencyVulns(ctx, repo.ID)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "ListDependencyVulns")
		return
	}
	apiVulns := make([]*api.DependencyVuln, len(vulns))
	pkgNamesMap := make(map[string]struct{})
	severityMap := make(map[string]struct{})
	targetMap := make(map[string]struct{})
	for i := range vulns {
		apiVulns[i] = convert.ToDependencyVuln(ctx, vulns[i])
		pkgNamesMap[vulns[i].PkgName] = struct{}{}
		severityMap[vulns[i].Severity] = struct{}{}
		targetMap[vulns[i].Target] = struct{}{}
	}
	pkgNames := make([]string, 0, len(pkgNamesMap))
	for pkg := range pkgNamesMap {
		pkgNames = append(pkgNames, pkg)
	}

	severities := make([]string, 0, len(severityMap))
	for severity := range severityMap {
		severities = append(severities, severity)
	}
	targets := make([]string, 0, len(targetMap))
	for target := range targetMap {
		targets = append(targets, target)
	}

	totalItems := len(apiVulns)
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
	ctx.Data["Pagination"] = map[string]interface{}{
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"HasPrev":     page > 1,
		"HasNext":     page < totalPages,
		"PrevPage":    page - 1,
		"NextPage":    page + 1,
		"Pages":       generatePages(totalPages),
	}
	start := (page - 1) * itemsPerPage
	end := start + itemsPerPage
	if end > totalItems {
		end = totalItems
	}
	paginatedList := apiVulns[start:end]
	ctx.Data["PageIsSecurity"] = true
	ctx.Data["PageIsDependencyScanning"] = true
	ctx.Data["DependencyVulns"] = paginatedList
	ctx.Data["PkgNames"] = pkgNames
	ctx.Data["Severities"] = severities
	ctx.Data["Targets"] = targets

	ctx.HTML(http.StatusOK, tplSecurityDependencyScanning)
}

// DepedencyVuln renders single depedency vulnerability page
func DepedencyVuln(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.security")
	ID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Vuln ID"})
		return
	}

	vulns, err := repo_model.GetDependencyVuln(ctx, ID)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "GetDependencyVuln")
		return
	}

	if vulns == nil {
		ctx.HTML(http.StatusNotFound, base.TplName("status/404"))
		return
	}

	apiVulns := make([]*api.DependencyVuln, 1)
	apiVulns[0] = convert.ToDependencyVuln(ctx, vulns)

	referenceList := strings.Split(apiVulns[0].References, ",")

	issueVulnList, err := issue_model.GetIssuesByVuln(ctx, apiVulns[0].ID)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "GetIssuesByVuln")
		return
	}

	var issueList issue_model.IssueList
	for _, issueVuln := range issueVulnList {
		issue, err := issue_model.GetIssueByIndex(ctx, issueVuln.RepoID, issueVuln.IssueIndex)
		if err != nil {
			log.Error(err.Error())
			return
		}
		issueList = append(issueList, issue)
	}

	ctx.Data["PageIsSecurity"] = true
	ctx.Data["PageIsDependencyScanning"] = true
	ctx.Data["DependencyVuln"] = apiVulns[0]
	ctx.Data["ReferencesList"] = referenceList
	ctx.Data["IssueList"] = issueList
	ctx.HTML(http.StatusOK, tplSecurityDependencyVuln)
}

// IacMisconfiguration renders single infrastructure as code misconfiguration page
func IacMisconfiguration(ctx *context.Context) {
	repo := ctx.Repo.Repository
	ctx.Data["Title"] = ctx.Tr("repo.security")

	iacFiltered := map[string]string{
		"type":        "",
		"location":    "",
		"severity":    "",
		"branch_name": "",
		"q":           "",
	}
	for k := range iacFiltered {
		if _type := ctx.FormTrim(k); _type != "" {
			iacFiltered[k] = _type
		}
	}
	iacMisconfigurationList, err := repo_model.ListRepoIacMisconfigurations(ctx, repo.ID, iacFiltered)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "ListRepoIacMisconfigurations")
		return
	}
	apiIacMisconfigurations := make([]*api.IacMisconfiguration, len(iacMisconfigurationList))
	for i := range iacMisconfigurationList {
		apiIacMisconfigurations[i] = convert.ToIacMisconfiguration(ctx, iacMisconfigurationList[i])
	}
	page := ctx.FormInt("page")
	if page < 1 {
		page = 1
	}
	itemsPerPage := 10
	totalItems := len(apiIacMisconfigurations)
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
	ctx.Data["Pagination"] = map[string]interface{}{
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"HasPrev":     page > 1,
		"HasNext":     page < totalPages,
		"PrevPage":    page - 1,
		"NextPage":    page + 1,
		"Pages":       generatePages(totalPages),
	}
	start := (page - 1) * itemsPerPage
	end := start + itemsPerPage
	if end > totalItems {
		end = totalItems
	}
	paginatedList := apiIacMisconfigurations[start:end]
	ctx.Data["PageIsSecurity"] = true
	ctx.Data["PageIsIaCMisconfiguration"] = true
	ctx.Data["IaCFilters"] = repo_model.GetListIaCFilter(ctx, repo.ID)
	ctx.Data["IaCFiltered"] = iacFiltered
	ctx.Data["IaCMisconfigurationList"] = paginatedList

	ctx.HTML(http.StatusOK, tplSecurityIaCMisconfig)
}

func parseCodeContent(codeContent string) []security_service.Line {
	var lines []security_service.Line
	for _, line := range strings.Split(codeContent, "\n") {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			lineNumberStr := strings.TrimSpace(parts[0])
			lineNumber, _ := strconv.Atoi(lineNumberStr)
			lineContent := line[len(lineNumberStr):]
			lines = append(lines, security_service.Line{
				Number:  lineNumber,
				Content: lineContent,
			})
		} else if len(parts) == 1 {
			lineNumberStr := strings.TrimSpace(parts[0])
			lineNumber, _ := strconv.Atoi(lineNumberStr)
			lines = append(lines, security_service.Line{
				Number:  lineNumber,
				Content: "",
			})
		}
	}
	return lines
}

// DetailIacMisconfiguration renders single detail iac misconfiguration page
func DetailIacMisconfiguration(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.security")
	ID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Iac Misconfiguration ID"})
		return
	}

	iacMisconfiguration, err := repo_model.GetRepoIacMisconfiguration(ctx, ID)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "GetRepoIacMisconfiguration")
		return
	}

	if iacMisconfiguration == nil {
		ctx.HTML(http.StatusNotFound, base.TplName("status/404"))
		return
	}

	apiIaCMisconfiguration := make([]*api.IacMisconfiguration, 1)
	apiIaCMisconfiguration[0] = convert.ToIacMisconfiguration(ctx, iacMisconfiguration)

	referenceList := strings.Split(apiIaCMisconfiguration[0].References, ",")

	ctx.Data["PageIsSecurity"] = true
	ctx.Data["PageIsIaCMisconfiguration"] = true
	ctx.Data["IaCMisconfiguration"] = apiIaCMisconfiguration[0]
	ctx.Data["ReferencesList"] = referenceList
	ctx.Data["CodeLines"] = parseCodeContent(apiIaCMisconfiguration[0].CodeContent)
	ctx.HTML(http.StatusOK, tplSecurityDetailIaCMisconfig)
}

// SecretDetection renders single secret detectuin page
func SecretDetection(ctx *context.Context) {
	repo := ctx.Repo.Repository
	ctx.Data["Title"] = ctx.Tr("repo.security")
	secretDetectionFiltered := map[string]string{
		"branch_name": "",
		"location":    "",
		"severity":    "",
		"q":           "",
	}
	for filterKey := range secretDetectionFiltered {
		if filterValue := ctx.FormTrim(filterKey); filterValue != "" {
			secretDetectionFiltered[filterKey] = filterValue
		}
	}
	secretDetectionList, err := repo_model.ListRepoSecretDetections(ctx, repo.ID, secretDetectionFiltered)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "ListRepoSecretDetections")
		return
	}
	apiSecretDetections := make([]*api.RepoSecretDetection, len(secretDetectionList))
	for i := range secretDetectionList {
		apiSecretDetections[i] = convert.ToSecretDetection(ctx, secretDetectionList[i])
	}
	page := ctx.FormInt("page")
	if page < 1 {
		page = 1
	}
	itemsPerPage := 10
	totalItems := len(apiSecretDetections)
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
	ctx.Data["Pagination"] = map[string]interface{}{
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"HasPrev":     page > 1,
		"HasNext":     page < totalPages,
		"PrevPage":    page - 1,
		"NextPage":    page + 1,
		"Pages":       generatePages(totalPages),
	}
	start := (page - 1) * itemsPerPage
	end := start + itemsPerPage
	if end > totalItems {
		end = totalItems
	}
	paginatedList := apiSecretDetections[start:end]

	ctx.Data["PageIsSecurity"] = true
	ctx.Data["PageIsSecretDetection"] = true
	ctx.Data["SecretDetectionList"] = paginatedList
	ctx.Data["SecretDetectionFilters"] = repo_model.GetListSecretDetectionFilter(ctx, repo.ID)
	ctx.Data["SecretDetectionFiltered"] = secretDetectionFiltered
	ctx.HTML(http.StatusOK, tplSecuritySecretDetection)
}

// DetailSecretDetection renders single detail secret detection page
func DetailSecretDetection(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.security")
	ID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Secret Detection ID"})
		return
	}

	secretDetection, err := repo_model.GetRepoSecretDetection(ctx, ID)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "GetRepoSecretDetection")
		return
	}

	if secretDetection == nil {
		ctx.HTML(http.StatusNotFound, base.TplName("status/404"))
		return
	}

	apiSecretDetection := make([]*api.RepoSecretDetection, 1)
	apiSecretDetection[0] = convert.ToSecretDetection(ctx, secretDetection)

	ctx.Data["PageIsSecurity"] = true
	ctx.Data["PageIsSecretDetection"] = true
	ctx.Data["SecretDetection"] = apiSecretDetection[0]
	ctx.Data["CodeLines"] = parseCodeContent(apiSecretDetection[0].CodeContent)
	ctx.HTML(http.StatusOK, tplSecurityDetailSecretDetection)
}

func ContainerScanning(ctx *context.Context) {
	repo := ctx.Repo.Repository
	ctx.Data["Title"] = ctx.Tr("repo.security")

	page := ctx.FormInt("page")
	if page < 1 {
		page = 1
	}
	itemsPerPage := 10

	containerScanningFiltered := map[string]string{
		"branch_name": "",
		"location":    "",
		"severity":    "",
		"q":           "",
	}
	for filterKey := range containerScanningFiltered {
		if filterValue := ctx.FormTrim(filterKey); filterValue != "" {
			containerScanningFiltered[filterKey] = filterValue
		}
	}
	containerScanningList, err := repo_model.ListRepoImageVulns(ctx, repo.ID, containerScanningFiltered)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "ListRepoContainerScannings")
		return
	}
	apiContainerScannings := make([]*api.ImageVuln, len(containerScanningList))
	for i := range containerScanningList {
		apiContainerScannings[i] = convert.ToImageVuln(ctx, containerScanningList[i])
	}
	totalItems := len(apiContainerScannings)
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
	ctx.Data["Pagination"] = map[string]interface{}{
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"HasPrev":     page > 1,
		"HasNext":     page < totalPages,
		"PrevPage":    page - 1,
		"NextPage":    page + 1,
		"Pages":       generatePages(totalPages),
	}
	start := (page - 1) * itemsPerPage
	end := start + itemsPerPage
	if end > totalItems {
		end = totalItems
	}
	paginatedList := apiContainerScannings[start:end]
	ctx.Data["PageIsSecurity"] = true
	ctx.Data["PageIsContainerScanning"] = true
	ctx.Data["ContainerScanningList"] = paginatedList
	ctx.Data["ContainerScanningFilters"] = repo_model.GetListImageVulnFilter(ctx, repo.ID)
	ctx.Data["ContainerScanningFiltered"] = containerScanningFiltered
	ctx.HTML(http.StatusOK, tplSecurityContainerScanning)
}

func generatePages(totalPages int) []int {
	const maxPagesToShow = 10
	var pages []int

	if totalPages <= maxPagesToShow {
		for i := 1; i <= totalPages; i++ {
			pages = append(pages, i)
		}
	} else {
		currentPage := 1
		start := max(1, currentPage-2)
		end := min(totalPages, start+maxPagesToShow-1)

		for i := start; i <= end; i++ {
			pages = append(pages, i)
		}
	}

	return pages
}

func DetailContainerScanning(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.security")
	ID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Container Scanning ID"})
		return
	}

	containerScanning, err := repo_model.GetRepoImageVuln(ctx, ID)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "GetRepoContainerScanning")
		return
	}

	if containerScanning == nil {
		ctx.HTML(http.StatusNotFound, base.TplName("status/404"))
		return
	}
	apiContainerScanning := make([]*api.ImageVuln, 1)
	apiContainerScanning[0] = convert.ToImageVuln(ctx, containerScanning)
	referenceList := strings.Split(apiContainerScanning[0].References, ",")

	ctx.Data["PageIsSecurity"] = true
	ctx.Data["PageIsContainerScanning"] = true
	ctx.Data["ContainerScanning"] = apiContainerScanning[0]
	ctx.Data["References"] = referenceList
	ctx.HTML(http.StatusOK, tplSecurityDetailContainerScanning)
}

func Dashboard(ctx *context.Context) {
	ctx.PageData["repoLink"] = ctx.Repo.RepoLink
	ctx.Data["PageIsSecurity"] = true
	ctx.Data["PageIsDashboard"] = true
	ctx.HTML(http.StatusOK, tplDashboard)
}

func VulnStatistic(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.security")
	repo := ctx.Repo.Repository

	results, err := repo_model.GetVulnScanStatisticsLast12Days(ctx, repo.ID)
	if err != nil {
		ctx.ServerError("GetVulnScanStatisticsLast12Days", err)
		return
	}

	chartData := make(map[string][]map[string]interface{})
	for scanType, quantities := range results {
		var dataPoints []map[string]interface{}
		for i, quantity := range quantities {
			date := time.Now().UTC().AddDate(0, 0, -11+i).Format("2006-01-02")
			dataPoints = append(dataPoints, map[string]interface{}{
				"date":   date,
				"amount": quantity,
			})
		}
		chartData[scanType] = dataPoints
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"chartData": chartData,
	})
}

func VulnStatisticWeekly(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.security")
	repo := ctx.Repo.Repository

	results, err := repo_model.GetVulnScanStatisticsLast12Weeks(ctx, repo.ID)
	if err != nil {
		ctx.ServerError("GetVulnScanStatisticsLast12Weeks", err)
		return
	}

	chartData := make(map[string][]map[string]interface{})
	for scanType, quantities := range results {
		var dataPoints []map[string]interface{}
		for i, quantity := range quantities {
			date := time.Now().UTC().AddDate(0, 0, -7*(11-i))
			weekStart := date.AddDate(0, 0, -int(date.Weekday())+1).Format("2006-01-02")

			dataPoints = append(dataPoints, map[string]interface{}{
				"date":   weekStart,
				"amount": quantity,
			})
		}
		chartData[scanType] = dataPoints
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"chartData": chartData,
	})
}

func VulnStatisticMonthly(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.security")
	repo := ctx.Repo.Repository

	results, err := repo_model.GetVulnScanStatisticsLast12Months(ctx, repo.ID)
	if err != nil {
		ctx.ServerError("GetVulnScanStatisticsLast12Months", err)
		return
	}

	chartData := make(map[string][]map[string]interface{})
	for scanType, quantities := range results {
		var dataPoints []map[string]interface{}
		for i, quantity := range quantities {
			date := time.Now().UTC().AddDate(0, -1*(11-i), 0)
			monthStart := date.Format("2006-01")

			dataPoints = append(dataPoints, map[string]interface{}{
				"date":   monthStart,
				"amount": quantity,
			})
		}
		chartData[scanType] = dataPoints
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"chartData": chartData,
	})
}

func SecurityVulnBySeverity(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.security")
	repo := ctx.Repo.Repository

	listDependencyVuln, err := repo_model.ListDependencyVulns(ctx, repo.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Fail to list Dependency Vulns"})
		return
	}
	iacFiltered := map[string]string{
		"type":     "",
		"location": "",
		"severity": "",
		"q":        "",
	}
	listIACVulns, err := repo_model.ListRepoIacMisconfigurations(ctx, repo.ID, iacFiltered)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Fail to list IAC Misconfiguration Vulns"})
		return
	}
	secretDetectionFiltered := map[string]string{
		"branch_name": "",
		"location":    "",
		"severity":    "",
		"q":           "",
	}
	listSecretVulns, err := repo_model.ListRepoSecretDetections(ctx, repo.ID, secretDetectionFiltered)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Fail to list Secret Vulns"})
		return
	}
	severityCount := map[string]int{
		"critical": 0,
		"high":     0,
		"medium":   0,
		"low":      0,
	}
	for _, vuln := range listDependencyVuln {
		switch vuln.Severity {
		case "CRITICAL":
			severityCount["critical"]++
		case "HIGH":
			severityCount["high"]++
		case "MEDIUM":
			severityCount["medium"]++
		case "LOW":
			severityCount["low"]++
		}
	}
	for _, vuln := range listIACVulns {
		switch vuln.Severity {
		case "CRITICAL":
			severityCount["critical"]++
		case "HIGH":
			severityCount["high"]++
		case "MEDIUM":
			severityCount["medium"]++
		case "LOW":
			severityCount["low"]++
		}
	}
	for _, vuln := range listSecretVulns {
		switch vuln.Severity {
		case "CRITICAL":
			severityCount["critical"]++
		case "HIGH":
			severityCount["high"]++
		case "MEDIUM":
			severityCount["medium"]++
		case "LOW":
			severityCount["low"]++
		}
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"severity_counts": severityCount,
	})
}

func SecurityVulnByStatus(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.security")
	repo := ctx.Repo.Repository

	listDependencyVuln, err := repo_model.ListDependencyVulns(ctx, repo.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Fail to list Dependency Vulns"})
		return
	}
	statusCount := map[string]int{
		"detected": 0,
		"resolved": 0,
		"reopen":   0,
	}
	for _, vuln := range listDependencyVuln {
		switch vuln.Label {
		case "Detected":
			statusCount["detected"]++
		case "Resolved":
			statusCount["resolved"]++
		case "Reopen":
			statusCount["reopen"]++
		}
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status_counts": statusCount,
	})
}

// CodeScanning renders single code scanning page
func CodeScanning(ctx *context.Context) {
	repo := ctx.Repo.Repository
	ctx.Data["Title"] = ctx.Tr("repo.security")

	page := ctx.FormInt("page")
	if page <= 1 {
		page = 1
	}

	codeVulnFiltered := map[string]string{
		"vuln_class":  "",
		"location":    "",
		"severity":    "",
		"branch_name": "",
		"q":           "",
	}
	for k := range codeVulnFiltered {
		if _type := ctx.FormTrim(k); _type != "" {
			codeVulnFiltered[k] = _type
		}
	}
	codeVulnList, err := repo_model.ListRepoCodeVulns(ctx, repo.ID, codeVulnFiltered)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "ListRepoCodeVulns")
		return
	}

	total := len(codeVulnList)

	codeVulnList, _, err = db.FindAndCount[repo_model.RepoCodeVuln](ctx, repo_model.RepoCodeVulnSearchOptions{
		RepoCodeVuln: repo_model.RepoCodeVuln{
			RepoID:     repo.ID,
			VulnClass:  codeVulnFiltered["vuln_class"],
			Target:     codeVulnFiltered["location"],
			Severity:   codeVulnFiltered["severity"],
			BranchName: codeVulnFiltered["branch_name"],
		},
		ListOptions: db.ListOptions{
			PageSize: setting.UI.IssuePagingNum,
			Page:     page,
		},
		Q: codeVulnFiltered["q"],
	})
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "PaginatedCodeVulnList")
		return
	}

	apiCodeVulns := make([]*api.RepoCodeVuln, len(codeVulnList))
	for i := range codeVulnList {
		apiCodeVulns[i] = convert.ToCodeVuln(ctx, codeVulnList[i])
	}

	pager := context.NewPagination(total, setting.UI.IssuePagingNum, page, 10)

	pager.AddParamString("q", codeVulnFiltered["q"])
	pager.AddParamString("vuln_class", codeVulnFiltered["vuln_class"])
	pager.AddParamString("severity", codeVulnFiltered["severity"])
	pager.AddParamString("location", codeVulnFiltered["location"])
	pager.AddParamString("branch_name", codeVulnFiltered["branch_name"])

	ctx.Data["Page"] = pager
	ctx.Data["PageIsSecurity"] = true
	ctx.Data["PageIsCodeScanning"] = true
	ctx.Data["PaginatedCodeVulnList"] = apiCodeVulns

	ctx.Data["CodeVulnFilters"] = repo_model.GetListCodeVulnFilter(ctx, repo.ID)
	ctx.Data["CodeVulnFiltered"] = codeVulnFiltered

	ctx.HTML(http.StatusOK, tplSecurityCodeScanning)
}

// CodeVuln renders single detail code vuln page
func CodeVuln(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.security")
	ID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Code Vuln ID"})
		return
	}

	codeVuln, err := repo_model.GetRepoCodeVuln(ctx, ID)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "GetRepoCodeVuln")
		return
	}

	if codeVuln == nil {
		ctx.HTML(http.StatusNotFound, base.TplName("status/404"))
		return
	}

	apiCodeVuln := make([]*api.RepoCodeVuln, 1)
	apiCodeVuln[0] = convert.ToCodeVuln(ctx, codeVuln)
	referenceList := strings.Split(apiCodeVuln[0].References, ",")

	ctx.Data["PageIsSecurity"] = true
	ctx.Data["CodeVuln"] = apiCodeVuln[0]
	ctx.Data["ReferencesList"] = referenceList
	ctx.HTML(http.StatusOK, tplSecurityCodeVuln)
}

// Enable or disable security scanning
func ActionRepoSecurityScanning(ctx *context.Context) {
	repo := ctx.Repo.Repository

	// Check if "dependency_scanning_branch_name" exists in the request form data
	if _, exists := ctx.Req.Form["dependency_scanning_branch_name"]; exists {
		DependencyScanningBranchName := ctx.FormString("dependency_scanning_branch_name")
		repo.DependencyScanningBranchName = DependencyScanningBranchName
	}

	// Check if "is_dependency_list_enabled" exists in the request form data
	if _, exists := ctx.Req.Form["is_dependency_list_enabled"]; exists {
		isDependencyListEnabled := ctx.FormBool("is_dependency_list_enabled")
		repo.IsDependencyListEnabled = isDependencyListEnabled
	}

	// Check if "iac_misconfiguration_branch_name" exists in the request form data
	if _, exists := ctx.Req.Form["iac_misconfiguration_branch_name"]; exists {
		IacMisconfigurationBranchName := ctx.FormString("iac_misconfiguration_branch_name")
		repo.IacMisconfigurationBranchName = IacMisconfigurationBranchName
	}

	// Check if "secret_detection_branch_name" exists in the request form data
	if _, exists := ctx.Req.Form["secret_detection_branch_name"]; exists {
		SecretDetectionBranchName := ctx.FormString("secret_detection_branch_name")
		repo.SecretDetectionBranchName = SecretDetectionBranchName
	}

	// Check if "code_scanning_branch_name" exists in the request form data
	if _, exists := ctx.Req.Form["code_scanning_branch_name"]; exists {
		CodeScanningBranchName := ctx.FormString("code_scanning_branch_name")
		repo.CodeScanningBranchName = CodeScanningBranchName
	}

	// Update the repository only with changed values
	if err := repo_service.UpdateRepository(ctx, repo, false); err != nil {
		ctx.Error(http.StatusInternalServerError, "UpdateRepository")
		return
	}

	permission, err := access_model.GetUserRepoPermission(ctx, repo, ctx.Doer)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "GetTeamRepos")
		return
	}

	ctx.JSON(http.StatusOK, convert.ToRepo(ctx, repo, permission))
}

func UpdateDependencyVulnLabel(ctx *context.Context) {
	ID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Vuln ID"})
		return
	}

	// Check if "label" exists in the request form data
	if _, exists := ctx.Req.Form["label"]; exists {
		label := ctx.FormString("label")
		err = issue_service.SyncIssueVulnStatus(ctx, ID, label, ctx.Doer)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Update vulnerability status fail"})
			return
		}

	}

	ctx.JSON(http.StatusOK, map[string]string{"message": "Update dependency vulnerability successfully"})
}

func UpdateCodeVulnLabel(ctx *context.Context) {
	ID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Vuln ID"})
		return
	}

	label := ctx.FormString("label")

	err = repo_model.UpdateCodeVulnLabel(ctx, ID, label)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Update vulnerability status fail"})
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{"message": "Update dependency vulnerability successfully"})
}
