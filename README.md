# Factory Fuse API

Hi! I'm Fuse, nice to meet you!
In this documentation, I will tell you how to deploy your own precious Factory Fuse instance.


# Deploy

So Factory Fuse needs to be deployed, to run properly, right?
So let's deploy it.

## Prerequisites

You need to have:

 - Git
 - Go compiler
 - Docker (and preferably, Docker Desktop)
 - The official PostgreSQL Docker image. (Will show how to gather it.)
 -  Some terminal knowledge
 
 Be sure that most of the programs are on modern versions. We don't want to deal with the antique, do we?

## Clone

When you got the Git, you need to setup your SSH key for authentication.
*Ssh!!* This project is *super secret*, we don't want anyone to snoop around our code.
(Not really, the code is public, but it's a procedure, I don't know why.)
Please [read here and follow the steps](https://docs.github.com/en/authentication/connecting-to-github-with-ssh/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent), before cloning code, otherwise you can view and download the code, but you cannot push new code.

If your ssh key is set up, let's download the code.
Enter your terminal this code:

    git clone https://github.com/factoryfuse/ffapi.git

You will recognize a new folder has been created called "ffapi" on your working directory. Find it and voila! The code is downloaded! Hurray!

## Setup DB

Now it's time to set the database.
Download the image for the PostgreSQL.
(This will take time depending on your internet connection.)

    docker pull postgres
  
Let's check if image has downloaded correctly.

    docker images

If you see `postgres` on the list, you are good to go!

Now, let's create a container for the database.

    docker run --name ffapi-db -e POSTGRES_PASSWORD=ffapi_pass -p 5432:5432 -v D:\Postgres_Data:/var/lib/postgresql/data -d postgres

    97f3983908c234eb3f45b0ceac7e367f0f0f7b8046b2906781e695dedbf0ef4e

It outputted something, *it's talking to me!*
Don't worry, it's just the container ID. Basically, it represents the container in the Docker.

Now, let's check if the container is running.

    docker ps

If you see `ffapi-db` on the list, well done!
Let's start the database.

    docker start ffapi-db 
    // You don't need to do that, mostly it's already running.

You have set up the database, now let's set up the code.

## Run

Now, you have the code and the database, let's run the code.
First, let's go to the code directory.

    cd ffapi

Now, let's run the code!

    $env:DBUSER = 'root'; // We need to set the environment variables.
    $env:DBPASS = 'ffapi_pass'; // Otherwise, the code won't be able to connect to the database.
    // The developer said that it's a temporary solution, and he will find a better way to do it.

    go run .

*Ooh, what's that? Fancy text again!*

If you don't see any text called, "ERROR", congrats, you have run the code!

## End

If you have any questions, just don't hesitate to message the lead developer, [Enes](https://github.com/eneshalat), for better brief. He said he will answer you if you are in the FactoryFuse DevTeam. Otherwise for the privacy of the project, he will not give internal details.

