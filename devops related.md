steps to create an ec2(elastic compute) instance in aws,

1) create instance and select everything mostly to default, for the network security, add 3 rules,

    a) port 22 for your own

    b) port 80 for http for everyone

    c) port 443 for https for everyone

2) save the private key which was downloaded earlier to the any path(c/users/shivang/.ssh/.pem file)

3) open the powershell and then go inside the directory where you saved the .pem file and run this command "ssh -i "magicmovies-server.pem" ec2-user@YOUR_PUBLIC_IP"

4) Few Linux commands,

    4.1) Navigation & Exploring,

        a) pwd = print working directory - prints the current path.

        b) ls = List - Lists all files and folders in your current directory.

        c) cd = Change Directory - Moves you into a different folder (like double-clicking a folder on Windows). Eg:) cd folder_name to go inside a folder. Type cd .. to go back up one level. Type just cd to jump straight back to your home folder.

    4.2) Creating & Managing Files,

        a) mkdir = Make Directory - Creates a brand new empty folder. Eg:) mkdir magicmovies creates a folder named magicmovies

        b) touch = Touch - Creates a blank, empty file. Eg:) touch index.html creates an empty file named index.html

        c) rm = Remove - Deletes files or folders permanently. Eg:) rm index.html deletes a file. To delete a folder and everything inside it, use rm -rf folder_name.

        d) cp = Copy - Copies a file from one place to another. Eg:) cp index.html backup.html makes a duplicate.

        e) mv = Move / Rename - Moves a file to a new folder, OR renames a file. Eg:)mv index.html home.html renames the file. mv home.html magicmovies/ moves it into the folder.

    4.3) Viewing & Searching inside Files,

        a) cat = Concatenate - Dumps the entire contents of a file directly onto your screen. Perfect for small files Eg:) cat home.html

        b) less = Less - Opens a file in a scrollable view. Perfect for huge log files. Eg:)less system.log. Use arrow keys to scroll, and press q to exit/quit.

        c) grep = Global Regular Expression Print - Searches for specific word or text pattern in files. Eg:) grep "err" sys.log will print only the lines containing the word "err".

        d) find = Find - Searches the entire system for files or folders by name. Eg:) find . -name "\*.html" searches your current folder (.) for any file ending in .html

    4.4) System Monitoring,

    	a) ps = Process Status - Shows snapshot of programs currently running. Eg:) Type ps aux to see detailed list of every single background process running on your server

    	b) top = Top - A live, dynamic dashboard showing CPU usage, RAM usage, and active tasks like Windows Task Manager. Eg:) Type top. Press q to close it.

    	c) df = Disk Free - Shows how much hard drive storage space you have left. Eg:) Type df -h (the -h makes it "human-readable" in Gigabytes instead of raw bytes).

    	d) free = Free Memory - Shows how much RAM memory is free and used. Eg:) Type free -m to see your RAM usage displayed cleanly in Megabytes.

    4.5)  Internet, Permissions & Admin Power,

    	a) ping = Ping - It will just check whether thr URL is reachable or not. Eg:) ping -c 4 google.com Linux will just send 4 packets and then stop otherwise it will keep on pinging forever.

    	b) curl = Client URL - Downloads data or webpages from the internet straight into your terminal. Eg:) curl https://google.com will download Google's HTML code. curl -I https://google.com , The capital -I tells curl to only grab the status headers instead of downloading the whole webpage code.

        c) chmod = Change Mode - Changes the security permissions of a file (who can read, write, or execute it). Eg:) chmod +x script.sh makes a script file executable so it can be run as a program

        d) sudo = SuperUser Do - Runs a command with absolute admin/root privileges like "Run as Administrator" in Windows. On a new EC2 instance, your regular user can't install software. You must use sudo. Eg:) sudo yum update or sudo apt update to update your server packages

        e) tulpn = tcp/udp/listening(active)/processes/numeric(80,443) - Used to tell information about the all the active ports and what processed is runningon them. Eg:) sudo ss -tulpn

    4.6) Services related,

    	a) systemctl = system control - manage Services. Eg:) sudo systemctl start nginx ,also sudo is necessary as linux treats all ports below 1024 as privileged(admin only)

    	b) start the service = sudo systemctl start SERVICE_NAME.     Eg:) **sudo systemctl start go-server**

    	c) stop the service =  sudo systemctl stop SERVICE_NAME.      Eg:) **sudo systemctl stop go-server**

    	d) restart the service = sudo systemctl restart SERVICE_NAME. Eg:) **sudo systemctl restart go-server**

    	e) status the service = sudo systemctl status SERVICE_NAME.   Eg:) **sudo systemctl status go-server**

    	f) enable service when ec2 instance boots = sudo systemctl enable SERVICE_NAME. Eg:) **sudo systemctl enable go-server**

    4.7) Logging related,

    	a) journalctl = Journal control - to see the logs for particular service. Eg:) **sudo journalctl -u go-server.** For live logs, **sudo journalctl -u go-server -f.** For recent logs, **sudo journalctl -u go-server --since "10 minutes ago"**

    4.8) Few other commands,

    	a) nano FILE_NAME = used to open the text-editor thing where you can write the code and hit "Ctrl+O" to save and "ctrl+X" to exit. Eg:) **nano main.go**

    	b) kill = stop process. Eg:) kill ANY_PID, **kill 33061**



5) There is a few difference between Linux process and services, process are workers that perform some task usually started by the user and dies as soon as they are finished, commands like ls, pwd, etc. Service is the special process that runs continuously in the background and started automatic as soon as ec2 boots up, like sshd, nginx, httpd, etc.

6) public ips are the one which anyone can access, these are the one through which computer connect to anyone on internet. Private ips are the one which are unique in the local network and can be same outside the local network, these are the one aws uses to communicate between ec2 instances in the same account.

7) difference between aws security groups, linux firewall and nginx,

    a) Security Groups - These are the virtual firewall implemented by the aws. It lives outside the ec2 instance, at aws network level. It inspects the traffic based on rules before 	traffic even touches the virtual machine/ec2 instance.

    b) Linux Firewall - These are the firewall implemented by the Linux OS. It lives inside the OS of ec2 instance.

    c) nginx - It is not a firewall at all, it is just an application. Specifically it is an web server or reverse proxy service. It lived inside the user space. It simply provides the 	resources based on the request.

8) Reverse Proxy - A reverse proxy is an application that sits between internet users and your backend web servers. When a user requests a website, Nginx receives the request first and silently forwards it to the correct internal application running on your server. The user never communicates with your backend code directly

    8.1) Reverse Proxy Working -

    	a) Problem without Reverse Proxy - This exposes your internal app port to the internet. It requires you to open that internal Port(5000) in your AWS Security Group.It looks ugly and unprofessional to users.

                        Browser
                           │
                           │ :8080
                           ▼
                    AWS Security Group
                           │
                           ▼
                          EC2
                           │
                           ▼
                    Go application
                           │
                           ▼
                         :8080

    	b) With Reverse Proxy - nginx running on port 80 receives the request. It reads the URL path and forwards the request to the mapped internal port where backend code is running and returns the data to Nginx and Nginx to user.

                        Browser
                           │
                           │ :80/443
                           ▼
                    AWS Security Group
                           │
                           ▼
                          EC2
                           │
                           ▼
                        Nginx :80
                           │
                           │ reverse proxy
                           ▼
                    Go service :8080

    8.2) Reverse Proxy advantages - Security, handling static data, SSL/HTTPS Encryption, routing traffic based on url to different ports

9) In order to install any packages in ec2 Amazon Linux 2023, run the following commands,

    1) Update your packages - **sudo dnf update -y**

    2) Install Go - **sudo dnf install golang -y** 

    3) Verify it works - **go version**

10) general flow of the request from browser to serve running on ec2 machine,

                    Your browser
                        │
                        │ :8080
                        ▼
                    Internet
                        │
                        ▼
                 AWS Security Group
                        │
                     allowed?
                        │
                        YES
                        │
                        ▼
                     EC2 Linux
                        │
                        ▼
                   Go process
                        │
                     listening
                       :8080
                        │
                        ▼
             "Hello from my EC2 Go server!"

11) In order to make the program able to run as service in Linux, do these steps. This is needed because when we run via "go run main.go", the program will stop as soon as we hit ctrl+c or closes the terminal.

    1) create an executable file for the code (go build -o go-server)

    2) Create a systemd service,

    	2.1) Create a service file - **sudo nano /etc/systemd/system/go-server.service**

    	2.2) Write service configurations. See below for example,

    		[Unit]
    		Description=Go Server
    		After=network.target

    		[Service]
    		User=ec2-user
    		WorkingDirectory=/home/ec2-user/go-server
    		ExecStart=/home/ec2-user/go-server/go-server
    		Restart=always
    		RestartSec=5

    		[Install]
    		WantedBy=multi-user.target

    3) Tell ec2 machine to look(allocates memory) for the newly added/updated background service files because it doesn't know till reloaded - **sudo systemctl daemon-reload**

    4) start the service - **sudo systemctl start go-server**

    5) check the status - **sudo systemctl status go-server**

12) Steps to configure Nginx in EC2 machine,
    1) Install Nginx - **sudo dnf install nginx -y**

    2) start and enable Nginx to run whenever EC2 boots - **sudo systemctl start nginx** and **sudo systemctl enable nginx**

    3) test the service - **curl http://localhost**

    4) open the Nginx config file - **sudo nano /etc/nginx/nginx.conf**

    5) locate the server section in the file and update the file like this, It will pass the incoming request on '/' to Go application listenting on port :8080,

        server {
            listen 80;
            server_name _;

            location / {
                proxy_pass http://127.0.0.1:8080;

                proxy_set_header Host $host;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto $scheme;
            }
        }

    6) Test the changes before reloading the service - **sudo nginx -t**, you should see syntax OK.

    7) reload the service - **sudo systemctl reload nginx**

13) In order to make the site secure (HTTPS), we use two things,
    13.1) Parties Involved,
        a) Lets Encrypt - It is aa open source, free certificate provider for the site.
        
        b) Certbot - It is an open source, free automated manager that talks with Lets Encrypt in order to obtain the certificate, updates the Nginx configurations and renews the certitficate automatically 30 days before certificate expiration.

    13.2) Complete flow of how HTTPS/Nginx/Browser works together,
        STEP 1) The user types https://example.com into their browser. The browser looks up your domain's IP address and sends a request to your server on Port 443

        STEP 2) Nginx intercepts the request on Port 443. Before exchanging any website data, Nginx and the user's browser perform a "handshake":
            2.1) Nginx presents the Let's Encrypt certificate to the browser.
            2.2) The browser verifies that the certificate is valid, hasn't expired, and actually matches your domain name.
            2.3) Once verified, they establish a secure, encrypted tunnel.

        STEP 3) Now that the connection is secure, Nginx decrypts the user's request. What Nginx does next depends on how your website is built:
            Scenario A) Static Website (HTML/CSS files) Nginx looks for the files on your server's hard drive and sends them back through the encrypted tunnel to the user. Port 80 is never touched.
            
            Scenario B) Web Application (Node.js, Python, Django, etc.) Most web apps run on an internal port (like localhost:3000 or localhost:8000). Nginx acts as a Reverse Proxy. It takes the decrypted request and forwards it to your app port internally. Your app processes it, sends the response back to Nginx, and Nginx encrypts it before sending it back to the user. 

        STEP 4) The user's browser receives the encrypted data, decrypts it using the keys established in Step 2, and displays your website with the padlock icon in    the address bar.

        NOTE) If user type the HTTP then request comes at port 80. Nginx automatically forwards the request to port 443 and send **status 301 moved permanently** to browser. Then evrything from STEP 1 starts.

    13.3) Configuring HTTPS on EC2 machine,
        1) Make sure port 443 is open in the security groups the EC2 machine.
        2) Intsall the Certbot - **sudo dnf install certbot python3-certbot-nginx -y**
        3) Allow Certbot to automatically request and install a security certificate for your specific domain using your existing Nginx web server configurations - **sudo certbot --nginx -d YOUR_DOMAIN**

                        Internet
                            │
                  ┌─────────┴─────────┐
                  │                   │
                HTTP                HTTPS
                 :80                 :443
                  │                   │
                  ▼                   ▼
             ┌───────────────────────────┐
             │           Nginx           │
             │                           │
             │   HTTP → HTTPS redirect   │
             │      TLS termination      │
             │       Reverse proxy       │
             └────────────┬──────────────┘
                          │
                          │ HTTP
                          ▼
                   ┌──────────────┐
                   │ Go service   │
                   │    :8080     │
                   └──────────────┘

14) Docker - It is an open source containerization platform used to store the code, runtime, dependencies, and libraries into one lightweight package called as container.

    14.1) Problem before Docker - 
        a) If applicaton 1 uses Go version 2 and application 2 uses Go version 3, then both applications can't run on same server, known as dependency hell problem.
        Workaround - Companies would often buy another expensive servers just to run single small application.
        b) Companies would often use VMs(Virtual Box/VMware/etc). A single pyhsical server would be split into multiple VMs that would run different applications, but these separate VMs has their own OS which would waste host system CPU/Memory/RAM/resources.

    14.2) Advantages of Docker - 
        a) Docker uses the resources of host OS only, thus making it extremely fast compared to traditional VMs.
        b) Two application with different dependencies can run on same server.
        c) Easy to scale up.
        d) Consolidated form of whole application with all its code, libraries and runtime, thus making it easy to run on different server without need to configure the server, docker will handle everything on its own.
        e) Consistency on all platform thus eliminating the "works on my computer" problem.

    14.3) Docker Componenets -
        a) Docker file (Recipe)    - A text file containing sequential instructions on how to assemble your application environment.
        b) Docker Image (Cake Mix) - A read-only, static snapshot built from the Dockerfile that contains the application code and libraries.
        c) Docker Container (Cake) - The live, running, and isolated instance created from a Docker Image. Multiple container can be made using single Image. **A container only stays alive while its main process is running.**
        d) Docker Engine (Oven) - underlying client-server application that builds, runs, and manages the containers.

    14.4) Configuring Docker on EC2 -
        a) Install Docker - **sudo dnf install docker -y**
        b) Start Docker - **sudo systemctl start docker**
        c) Start docker whenever EC2 boots - **sudo systemctl enable docker**
        d) Install the hello-world image - **sudo docker run hello-world**
        e) See all the container - **sudo docker ps -a**, ps - process status, a - all(running or stopped)
        f) See all docker images - **sudo docker images**

    14.5) Dockerize the application(Go) -
        a) Travel into the application directory, cd DIRECTORY_NAME
        b) Create the dockerfile, nano dockerfile
        c) Write the code in the file and step-by-step breakdown of what every single line does:

            # It tells Docker to pull an official pre-built image that contains Go version 1.24 running on top of Alpine Linux. Alpine Linux is the lightweight Linux distribution. Using this means your container automatically has the Go compiler and tools ready to go, without bringing in unnecessary, heavy operating system bloated files. 
            FROM golang:1.26-alpine

            # Creates and switches to a working directory inside the container like (cd /app) or (mkdir /app)
            WORKDIR /app

            # Copies files from your local computer into the container's working directory (/app).
            COPY go.mod ./
            COPY main.go ./

            # Compiles your Go source code into a runnable executable file and named it go-server.
            RUN go build -o go-server .

            # Documents which network port the application inside the container is listening on. This line doesn't actually open the port to the public internet yet.
            EXPOSE 8080

            # Sets the Main Process that runs when the container boots up. When you type docker run, Docker launches the compiled ./go-server binary. As long as your Go code has that infinite web-server loop listening for traffic, this process stays active, and your container stays alive.
            CMD ["./go-server"]

        d) Create the Image from docker file,
            **sudo docker build -t go-server .**, -t - tags the image with the name we provided "go-server" otherwise docker will name it itself in hexadecimal form which is tough to remember, . - It tells that the dockerfile and all the required files are present in current directory where this command is being executed.

        e) Create the docker Container from the Image,
            **
            sudo docker run -d \
            --name go-server \
            -p 8080:8080 \
            go-server
            **

            Explanation,
                1) docker run - The core command that tells Docker to create a brand new container instance and start it up.
                2) d - run in background, detached mode otherwise the terminal would show the logs and everything.
                3) name go-server - Assigns a specific, human-readable name (go-server) to this running container.
                4) p 8080:8080 - (EC2 port -> Mapping -> Container) This links a network port on your host machine to a docker which uses the helper process called docker-proxy(route mapping between EC2 instance port and internal containers). This forwards the traffic to container applications. In order to see the port mapping, **sudo docker port go-server**
                5) go-server - This tells Docker which image to use as the template to build this container.
            
        f) Now at this moment the container should be up and running for viewing the application logs - **sudo docker logs go-server**
    
    14.6) Multi stage build process - This is a docker technique where we use multiple "from" in a single docker file to create multiple temporary stages, where 
        result of previous stages can be used in final production stage.
        Advantages,
            a) Tiny Image Size: You leave behind heavy development tools (like the Go compiler) and raw source code. Your final image only contains the finished application executable. A typical Go app image drops from 800MB+ down to just 15MB.
            b) More Security: Hackers often try to neutralize the compiler and other tools in order to exploit the container, removing these from production image increases security.
            c) No messy cleanup

    14.7) Process of moving from single stage build to multi stage build,
        a) go into the application directory.
        b) go to the Dockerfile, **nano Dockerfile**
        c) Write this code,

            # Stage 1: Build, does all the heavy lifting and creates the executable binary of go application.
            FROM golang:1.26-alpine AS builder

            WORKDIR /app

            COPY go.mod ./
            COPY main.go ./

            RUN go build -o go-server .

            # Stage 2: Runtime, creates the new slate and uses only the necessary linux and executable from stage 1 for the final image.
            FROM alpine:latest

            WORKDIR /app

            COPY --from=builder /app/go-server .

            EXPOSE 8080

            CMD ["./go-server"] 

        d) stop and delete the already running container, **sudo docker stop go-server** and **sudo docker rm go-server**
        e) build image, **sudo docker build -t go-server:v2 .**, run the image, **sudo docker run -d \--name go-server \-p 8080:8080 \go-server:v2**

    14.8) Deleing the container, **sudo docker rm CONTAINER_NAME** does not delete the image. In order to delete image, **sudo docker rmi IMAGE_NAME**

    14.9) Docker Networking, 
        a) Docker networking allows the isolated containers to talk to each other, to the EC2 machine and over internet.
        b) Advantage of docker networking,
            1) We can run the sensitive application like database without exposing them on the EC2 machine thus increasing security.
            2) When containers are placed on same network, docker starts acting as a DNS, resolving container IP address with their name automatically.
        c) Steps to make containers communicate,
            1) Create a custom network, **sudo docker network create my-app-network**
            2) Start Container 1 ( Eg: Database),       **sudo docker run -d --name db-service --network my-app-network redis**
            3) Start Container 2 ( Eg: Go Web Server),  **sudo docker run -d --name go-server -p 8080:8080 --network my-app-network go-server**
            4) We can connect with database container from Go container like, db-service:6379, where db-service is the container name and port on which it is running.
        d) In order to hide the application from exposing to the EC2 machine, omit the "-p" flag in command. 
        e) Each docker container has its own separate IP address, network, ports, etc.
                ┌─────────────────── EC2 Host Machine ───────────────────┐
                │                                                        │
                │  ┌─────────────────┐             ┌─────────────────┐   │
                │  │   Container 1   │             │   Container 2   │   │
                │  │                 │             │                 │   │
                │  │ Go App (p:8080) │             │ Go App (p:8080) │   │
                │  │ IP: 172.17.0.2  │             │ IP: 172.17.0.3  │   │
                │  │                 │             │                 │   │
                │  │ localhost =     │             │ localhost =     │   │
                │  │ Container 1     │             │ Container 2     │   │
                │  └────────┬────────┘             └────────┬────────┘   │
                │           │                               │            │
                │           └───────────────┬───────────────┘            │
                │                           │                            │
                │                Docker Virtual Bridge (`docker0`)       │
                └───────────────────────────┼────────────────────────────┘
                                            │
                                Port Mapping (-p 8080:8080)
                                            │
                                Outside World / Internet
                            
    14.10) Docker Compose,
        a) Docker compose is a tool that allows to manage and run multi container application using single configuration file.
        b) Advantages,
            1) Instead of typing each complex, lengthy command manually, we can use single file to handle everything.
            2) If we want to disable all the containers or system cleanup, instead of removing each container manually, we can simply disable the compose file.
            3) Compose file automatically creates the private network for us.
        c) Process of using Docker Compose(yaml) file,
            1) Go into the application directory, **cd go-server**
            2) Create the compose file, **nano compose.yaml**
            3) Inside that write,

                services:
                    go-server:
                        build: .
                        container_name: go-server
                        ports:
                            - "8080:8080"

                here, we write the services in the compose file. "go-server" is our service name which will build the image by taking the docker file present in this directory. It will automatically creates the container named as "go-server" which is mapped on port 8080.
            4) stop and delete the already running container, **sudo docker stop go-server** and **sudo docker rm go-server**
            5) Run the compose file, **sudo docker compose up -d**
        d) Docker compose commands,
            1) to start the service, **sudo docker compose up -d**
            2) see containers, **sudo docker compose ps**
            3) stop the service, **sudo docker compose stop**
            4) stop and remove the service/containers, **sudo docker compose down**
            5) check for the changes in files and accordingly start the container, **sudo docker compose up -d --build**

    14.11) Volumes in Docker,
        a) Volumes are the data which is stored outside the container. Lets say there is a database and for some reason if that container is down and deleted then all the database files will be deleted as well, so we store sensitive data outside the container in order to use it again.
        b) Process of using Volumes,
            1) do, **nano compose.yaml**
            2) Inside that write,

                services:
                    go-server:
                        build: .
                        container_name: go-server
                        ports:
                            - "8080:8080"
                        depends_on:
                            - mongo

                    mongo:
                        image: mongo:8
                        container_name: mongo
                        volumes:
                            - mongo-data:/data/db

                volumes:
                    mongo-data:

        c) Docker compose will perform these tasks,
                - Create go-server image
                - Create mongo image
                - Create both containers
                - Create private network
                - Create Volumes
                - start the containers
        d) In order to test test the mongo reachibility from go, do,
            1) Enter Go Container(container_name from compose file), **sudo docker exec -it go-server sh**
            2) ping mongo, **ping mongo**
            3) type "exit" for coming back to ec2 machine.
        e) In order to test the volumes persistence, do,
            1) check the volumes, **sudo docker volume ls**
            2) Delete the compose service, **sudo docker compose down**
            3) check for the containers. It will not show any containers, **sudo docker ps**
            4) check for the volumes. It will show same volume as before deleting the compose, **sudo docker volume ls**
            5) start everything again, **sudo docker compose up -d**. It will create the containers but volumes will be reused.
            6) In order to delete the volumes as well, **sudo docker compose down -v**         

    14.12) Environment variables in docker compose,
        a) Environement variables are the secrets used by the application, they can differ based on the environments(localhost, dev, prod).
        b) Adding them in compose file, **nano compose.yaml**

                go-server:
                    build: .
                    container_name: go-server
                    ports:
                        - "8080:8080"

                    environment:
                        MONGO_URI: mongodb://mongo:27017
                        REDIS_ADDR: redis:6379

                    depends_on:
                        - mongo
                        - redis

        c) Some sesnsitive secrets should not be stored in compose file since it will be pushed on github and secrest can be exposed. In order to handle this, create the .env file and use that.
        d) Steps to use the sensitive secrets safely are,
            1) go inside the app directory, **cd go-server**
            2) create the env file, **nano .env** and write the secrets there
            3) update the compose file,

                    go-server:
                        build: .
                        container_name: go-server
                        ports:
                            - "8080:8080"

                        env_file:
                            - .env

                        depends_on:
                            - mongo
                            - redis

            4) Add .env to .gitignore, **nano .gitignore**

        e) Never put the secrets in the image itself because it will become the part of that image, instead we keep the image generic and .env file specified to the environments(dev, test, prod, etc).
        f) Putting secrets in the .env file does not make it safe automatically, anyone with the server files access can view that. So generally we use the amazon secret manager.
        g) In order to view the environments, **sudo docker exec SERVICE_NAME env**

    14.13) Health Checks in docker compose,
        a) We use "depends_on" settings in compose file, 
           Eg:)depends_on:
                - mongo
                - redis  
           it is often misunderstood that first mongo and redis gets ready then go server will start. But actually it decides the startup order, start mongo --> start redis --> start go. It is possible that go server starts even when mongo is still initializing. 
        b) health check actually means that run particular command periodically and tell if container is actually healty.
        c) process of adding health check in compose file, **nano compose.yaml**
            1)  update depends on and healthcheck in compose.yaml,
                    go-server:
                        build: .
                        container_name: go-server
                        ports:
                            - "8080:8080"

                        environment:
                            MONGO_URI: mongodb://mongo:27017
                            REDIS_ADDR: redis:6379

                        depends_on:
                            mongo:
                                condition: service_healthy
                            redis:
                                condition: service_healthy

                    redis:
                        image: redis:8
                        container_name: redis

                        healthcheck:
                            test: ["CMD", "redis-cli", "ping"]
                            interval: 10s
                            timeout: 5s
                            retries: 5
                    
            2)  Understand each part,
                    test - It asks redis if it is responding
                    interval - ask redis if it is responding after every 10 seconds
                    timeout - if redis won't respons, wait for 5 second before marking the check failed
                    retries - retry 5 times and after 5 consecutive fails, mark the container unhealthy
        d) Now the depends_on actually becomes more fail safe as go server will only start when other dependencies are healthy.

                          Compose starts
                                │
                        ┌───────┴────────┐
                        ↓                ↓
                      Mongo            Redis
                        │                │
                    healthcheck      healthcheck
                        │                │
                     healthy?         healthy?
                        │                │
                        └───────┬────────┘
                                ↓
                              Go API

        e) to check the status of any service, **sudo docker inspect SERVICE_NAME --format='{{.State.Health.Status}}'**
        f) restarting of container after failure, docker give us few restart policies,
            1) no             = Don't restart
            2) on-failure     =	Restart when process exits with an error
            3) always	      = Always restart
            4) unless-stopped = Restart unless you explicitly stopped it

            Example) services:
                        go-server:
                            build: .
                            container_name: go-server
                            ports:
                                - "8080:8080"

                            restart: unless-stopped

        


15) CI/CD,

Experience with AWS with knowledge of AWS Services like Autoscaling, ELB, ElastiCache, SQS, SNS, RDS, S3, Serverless Architecture, Lambda, Gateway, and Amazon DynamoDB, etc... or similar technology stack
