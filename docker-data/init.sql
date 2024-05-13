-- MySQL dump 10.13  Distrib 8.0.36, for Linux (x86_64)
--
-- Host: localhost    Database: sql_lab_cms
-- ------------------------------------------------------
-- Server version	8.0.36-0ubuntu0.23.10.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `contest_issues`
--

USE sql_lab_cms;
DROP TABLE IF EXISTS `contest_issues`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `contest_issues` (
  `id` int NOT NULL AUTO_INCREMENT,
  `contestId` int NOT NULL,
  `issueId` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `contestId` (`contestId`),
  KEY `issueId` (`issueId`),
  CONSTRAINT `contest_issues_ibfk_1` FOREIGN KEY (`contestId`) REFERENCES `contests` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `contest_issues_ibfk_2` FOREIGN KEY (`issueId`) REFERENCES `issues` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `contest_issues`
--

LOCK TABLES `contest_issues` WRITE;
/*!40000 ALTER TABLE `contest_issues` DISABLE KEYS */;
INSERT INTO `contest_issues` VALUES (1,1,1),(2,1,2),(3,2,1),(4,3,1),(5,3,3),(6,4,4),(7,4,5),(8,5,2),(9,5,3),(19,1,5);
/*!40000 ALTER TABLE `contest_issues` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `contest_users`
--

DROP TABLE IF EXISTS `contest_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `contest_users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `contestId` int NOT NULL,
  `userId` int NOT NULL,
  `participantRole` tinyint NOT NULL,
  `startTime` datetime DEFAULT NULL,
  `endTime` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `contestId` (`contestId`),
  KEY `userId` (`userId`),
  CONSTRAINT `contest_users_ibfk_1` FOREIGN KEY (`contestId`) REFERENCES `contests` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `contest_users_ibfk_2` FOREIGN KEY (`userId`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `contest_users`
--

LOCK TABLES `contest_users` WRITE;
/*!40000 ALTER TABLE `contest_users` DISABLE KEYS */;
INSERT INTO `contest_users` VALUES (1,1,1,1,NULL,NULL),(2,1,2,1,NULL,NULL),(3,2,1,1,NULL,NULL),(4,3,1,1,NULL,NULL);
/*!40000 ALTER TABLE `contest_users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `contests`
--

DROP TABLE IF EXISTS `contests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `contests` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(255) NOT NULL,
  `title` varchar(255) DEFAULT NULL,
  `openTime` datetime DEFAULT NULL,
  `closeTime` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `contests`
--

LOCK TABLES `contests` WRITE;
/*!40000 ALTER TABLE `contests` DISABLE KEYS */;
INSERT INTO `contests` VALUES (1,'contest1','Programming Contest','2024-04-15 08:00:00','2024-04-20 12:00:00'),(2,'contest2','Math Contest','2024-05-01 09:00:00','2024-05-01 11:00:00'),(3,'contest3','Science Quiz','2024-05-10 14:00:00','2024-05-10 16:00:00'),(4,'contest4','History Bee','2024-05-15 10:00:00','2024-05-15 12:00:00'),(5,'contest5','Art Competition','2024-05-20 13:00:00','2024-05-20 15:00:00');
/*!40000 ALTER TABLE `contests` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `issues`
--

DROP TABLE IF EXISTS `issues`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `issues` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(255) NOT NULL,
  `title` varchar(255) DEFAULT NULL,
  `questionContent` text,
  `point` int DEFAULT NULL,
  `limitedTime` float DEFAULT NULL,
  `useTables` varchar(255) DEFAULT NULL,
  `executeType` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `issues`
--

LOCK TABLES `issues` WRITE;
/*!40000 ALTER TABLE `issues` DISABLE KEYS */;
INSERT INTO `issues` VALUES (1,'SQL1','Customer1','<p>Viết câu truy vấn trả về tất cả thông tin của 3 khách hàng đầu tiên trong bảng <b>Customer</b></p>',10,1,'Customer','SELECT'),(2,'SQL2','Category insert basic','Query the NAME field for all American cities in the CITY table with populations larger than 120000. The CountryCode for America is USA.\n\nThe CITY table is described as follows:\nhttps://s3.amazonaws.com/hr-challenge-images/8137/1449729804-f21d187d0f-CITY.jpg',15,1,'Category','INSERT'),(3,'SQL3','Category update basic','<p>Cập nhật lại description = categoryName + categoryId đối với các Category có categoryId lớn hơn 4</p>',20,1,'Category','DELETE'),(4,'SQL4','create table Product','<p>This is the question content for Problem 4.</p>',12,1,'table4','CREATE_TABLE'),(5,'SQL5','Problem 5','<p>This is the question content for Problem 5.</p>',18,1,'table5','CREATE_TABLE'),(13,'SQL6','TEST TABLE 2','<p>This is the question content for Problem 5.</p>',NULL,200,'Customer,Category','SELECT');
/*!40000 ALTER TABLE `issues` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `submitions`
--

DROP TABLE IF EXISTS `submitions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `submitions` (
  `id` int NOT NULL AUTO_INCREMENT,
  `contestUserId` int NOT NULL,
  `contestIssueId` int NOT NULL,
  `status` varchar(255) DEFAULT NULL,
  `compiler` varchar(255) NOT NULL,
  `srcCode` text,
  `submitTime` datetime DEFAULT NULL,
  `executionTime` float DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `contestUserId` (`contestUserId`),
  KEY `contestIssueId` (`contestIssueId`),
  CONSTRAINT `submitions_ibfk_1` FOREIGN KEY (`contestUserId`) REFERENCES `contest_users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `submitions_ibfk_2` FOREIGN KEY (`contestIssueId`) REFERENCES `contest_issues` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `submitions`
--

LOCK TABLES `submitions` WRITE;
/*!40000 ALTER TABLE `submitions` DISABLE KEYS */;
INSERT INTO `submitions` VALUES (1,1,1,'AC','mysql',NULL,'2024-04-23 23:24:12',3),(2,1,1,'AC','mysql',NULL,'2024-04-23 23:24:12',2),(3,1,2,'WA','mysql',NULL,'2024-04-23 23:24:12',4),(4,2,1,'TLE','sql server',NULL,'2024-04-23 23:24:12',NULL);
/*!40000 ALTER TABLE `submitions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `testcases`
--

DROP TABLE IF EXISTS `testcases`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `testcases` (
  `id` int NOT NULL AUTO_INCREMENT,
  `issueId` int DEFAULT NULL,
  `outputPath` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `issueId` (`issueId`),
  CONSTRAINT `testcases_ibfk_1` FOREIGN KEY (`issueId`) REFERENCES `issues` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `testcases`
--

LOCK TABLES `testcases` WRITE;
/*!40000 ALTER TABLE `testcases` DISABLE KEYS */;
INSERT INTO `testcases` VALUES (1,1,'issue_1.json'),(2,2,'issue_2.json'),(3,3,'issue_3.json'),(4,3,'issue_4.json'),(9,13,'TEST_TABLE_21713973188518.json');
/*!40000 ALTER TABLE `testcases` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `firstName` varchar(255) DEFAULT NULL,
  `lastName` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'v','v','Vinh','Bui'),(2,'g','v','Giang','Bui');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-05-10 16:24:34