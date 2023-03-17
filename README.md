# Event Authenticator

## Install

```bash
    go get golang.org/x/oauth2
    go get github.com/asaskevich/govalidator
```

## Usage

```go
    func main() {
	flag.Parse()
	fmt.Println(mode)
	defer utilities.TimeTrack(time.Now(), fmt.Sprintf("Wedding API Service"))
	defer func() {
		fmt.Print("ef")
		if e := recover(); e != nil {
			log.Panicln(e)
			main()
		}
	}()

	//load env
	var config Config
	err := utilities.LoadEnvFromFile(&config, configPrefix, configSource)
	if err != nil {
		log.Fatalln(err)
	}

	//load container
	container, err = NewContainer(config)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Server is running at : " + config.Binding)
	http.ListenAndServe(config.Binding, NewAPIv1(container))
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&configPrefix, "configPrefix", "wedding", "config prefix")
	flag.StringVar(&configSource, "configSource", ".env", "config source")

}
```

## Middleware

## Details

## API

###LOGIN

- `/auth/sign-in` sign-in
- `/auth/sign-up` sign-up
- `/auth/verify:headers`

###USER
\*GET

- `/auth/user/` get all user
- `/auth/user/:id` get a user with id

\*POST

- `/auth/user/` add a user

\*PUT

- `/auth/user/update/:id` update a user

\*DELETE 

-`/auth/user/deletebyid/:id` delete a user by id
- `/auth/user/deletebyphone/:phonenumber` delete a user by phone number

###LOGIN WITH FACEBOOK
\*POST

- `/auth/facebook/sign-in`

###LOGIN WITH GOOGLE
\*POST

- `auth/google/sign-in`

## Verified authentication

## Register oauth2

\*Google Auth

1. Create a new project: https://console.developers.google.com/project
2. Choose the new project from the top right project dropdown (only if another project is selected)
3. In the project Dashboard center pane, choose "API Manager"
4. In the left Nav pane, choose "Credentials"
5. In the center pane, choose "OAuth consent screen" tab. Fill in "Product name shown to users" and hit save.
6. In the center pane, choose "Credentials" tab.
   - Open the "New credentials" drop down
   - Choose "OAuth client ID"
   - Choose "Web application"
   - Application name is freeform, choose something appropriate
   - Authorized origins is your domain ex: https://example.mysite.com
   - Authorized redirect URIs is the location of oauth2/callback constructed as domain + /auth/google/callback, ex: https://example.mysite.com/auth/google/callback
   - Choose "Create"
7. Take note of the Client ID and Client Secret

\*Facebook Auth

1. From https://developers.facebook.com select "My Apps" / "Add a new App"
2. Set "Display Name" and "Contact email"
3. Choose "Facebook Login" and then "Web"
4. Set "Site URL" to your domain, ex: https://example.mysite.com
5. Under "Facebook login" / "Settings" fill "Valid OAuth redirect URIs" with your callback url constructed as domain + /auth/facebook/callback
6. Select "App Review" and turn public flag on. This step may ask you to provide a link to your privacy policy.
