###

# Authenication vs Authorization 1 service

# login facebook

- Access to URL : https://developers.facebook.com/
- Create Account
- Create App with Website domain, your app have ClientID(App ID) and CLient Secret(Secret)
- In Setting/Basic : Fill Privacy Policy URL <i style="color: red"> Not Set as Localhost... </i> and then Turn ON the Development Mode

- Documentation: https://www.npmjs.com/package/react-facebook
- Using Login-Button to create button on interface : <i>Input your app ID</i>
- The http in interface must be secure : HTTPS - Because Enforce HTTPS in your Facebook App is ON
- Run app and click the button
- in tab console will return accessToken in object

Service

- GET : https://graph.facebook.com/me?fields=" + your permission + "&access_token=" + accessToken. Web will return object include all fields which you request
- With this object you can use it as you like
