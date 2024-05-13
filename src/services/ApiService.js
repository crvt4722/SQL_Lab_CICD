import axios from 'axios';
class ApiServiceClass {
  validToken() {
    return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwic3ViIjoiYWRtaW4iLCJyb2xlIjoxLCJleHAiOjE3NzA5NDgxNDd9.DX9Cs1xQ44JBnPHYr0gm_BMHqnPEyMK1WsE-roVYcCg";
  }
  async get(url) {
    let config = {
      headers: {
        'accept': 'application/json',
        'Authorization': 'Bearer ' + this.validToken()
      }
    };
    const response = await axios.get(
      url,
      config
    );
    return response;
  }

  async put(url, data) {
    let config = {
      headers: {
        'accept': 'application/json',
        'Authorization': 'Bearer ' + this.validToken()
      }
    };
    const response = await axios.put(
      url,
      data,
      config
    );
    return response;
  }

  async delete(url) {
    let config = {
      headers: {
        'accept': 'application/json',
        'Authorization': 'Bearer ' + this.validToken()
      }
    };
    const response = await axios.delete(
      url,
      config
    );
    return response;
  }

  async post(url, data) {
    let config = {
      headers: {
        'accept': 'application/json',
        'Authorization': 'Bearer ' + this.validToken()
      }
    };
    const response = await axios.post(
      url,
      data,
      config
    );
    return response;
  }
}

const ApiService = new ApiServiceClass;
export default ApiService;

