import axios from "axios";

class Auth {

  // Attempt to log a user in
  login(user) {
    return axios
      .post('/user/login', user).then(response => {

        // Save jwt
        if(response.data.token) {
          localStorage.setItem('user', JSON.stringify(response.data));
        }

        return response.data
      }, (error) => {
        console.log(error)
      });
  }

  // Logout a user
  logout() {
    localStorage.removeItem('user');
  }

  // Create a new user
  signup(user) {
    axios.post('/user/signup', user).then(response => {
      if(response.code == 200) {
        return true;
      } else {
        return false;
      }
    }, (error) => {
      console.log(error)
      return false
    });
  }

  // Validate token
  valitity() {
     return axios.get('/token/validate',  {headers: authHeader()})
      .then(response => response.data)
      .catch(err => console.error(err))
  }
}

export default new Auth();

// set auth header for API
export function authHeader() {
  let user = JSON.parse(localStorage.getItem('user'))

  if (user && user.token) {
    return { Authorization: 'Bearer ' + user.token };
  } else {
    return {};
  }
}

// Guard pages for authenticated users only
export const authGuard = (to, from, next) => {
  let auth = new Auth();

  // Test valitity of current user
  auth.valitity().then(response => {
    if (response.valid) {
      // Allow
      return next();
    } else {
      // Otherwise, log in
      window.location.replace("/login");
    }
  })  
};