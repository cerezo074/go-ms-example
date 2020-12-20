1. Install Docker from their website

2. Download postgres docker image with the following command:
    - docker pull postgres:13-alpine 
    
For this version we use 13-alpine which is a small an lighweight postgress version. Remenber to install the lastest alpine version.

3. Create postgres container with the following command:
    - docker run --name postgres13 -p 5432:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=password -d postgres:13-alpine

The container name is postgres13 and the image name is postgres:13-alpine. This commnad will let us or container created and running.
Notice that we create an user with its password, they are *admin* and *password* respectively.

4. Install golang migration tool(golang-migrate), from OSX this can be done easily with homebrew through the following command:
    - brew install golang-migrate

5. Create our database with this command:
    - docker exec -it postgres13 createdb --username=admin --owner=admin freefortalking 

6. Run migrations with the following command:
    - migrate -path db/migration -database "postgres://admin:password@localhost:5432/freefortalking?sslmode=disable" -verbose up

Before runnig our migrations located in db/migrations folder, please check that our postgres container is up(status) and running with the following command.
- docker ps -a

Note: almost all commands are located on Makefile, so you dont need to copy them just need to check which command is that you need and run:
    - make <command>