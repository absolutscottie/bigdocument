# bigdocument

Table of Contents
=================
* [Description of the Problem](#description-of-the-problem)
* [Approach](#approach)
* [Project Layout](#project-layout)
* [Design](#design)
* [Required packages](#required-packages)
  * [Logrus](#logrus)
  * [Mux](#mux)
  * [Mongo DB Driver](#mongo-db-driver)
* [Running The Server](#running-the-server)
* [Executing Tests](#executing-tests)

### Description of the problem
For the purposes of this challenge, assume the words can be arbitrarily long UTF8 strings.

Requirements:
1) Implement "PUT" HTTP endpoint which will save/replace file on the disk.
2) Implement "GET "HTTP endpoint which will return a "processed" file without word duplicates.
3) Order of the input/output can be ignored.
4) Assume you cannot hold the entire file in memory.

Priority: Normal
Estimation: __SP

Acceptance criteria:
We will look into:
Code quality
Performance, complexity
There are no limitations to tools/programming language usage.
You're free to use any framework, but we'd like to see a short explanation of your decision. It will be good to have a short note on why the specific architecture approach was chosen.
The working code is covered with tests code. There is no need to cover the codebase for 100% with tests, but the key parts should be covered so we can run them. Also, instrumented tests can be skipped.
The solution should be shared using any VCS hosting.

>Note that I asked a clarification question regarding the size of any given word in the file. I am able to expect that each word would be no larger than 128 characters.

>Please also note that I would likely split this into a couple of different stories rather than roll it into one so that work could be done in parallel. At a high level, 3 x 3 point stories seems appropriate: one for the ingest work, one for the egress work and one to do the data storage part.

### Approach
My approach to the problem is that I will solve the problem of word repeats at ingest by using a data structure where I can prevent duplicates from being stored. The API does not require that the file ever needs to be returned in it's vanilla state, so there's no need to store it unaltered. There were multiple wasy to accomplish this from a database point of view, so the solution I chose was this:
* Parse the input as it is being consumed (by line breaks)
* Store the words in named tables or collections
* Insert each word to the database separately 
* Rely on a unique index/key to prevent duplicates

I could have accomplished the above with most databases that I'm familiar with but I chose Mongo DB for two reasons:
* The internet told me that it handles large amounts of data well
* Science. I wanted to get some experience with Mongo DB and this was a good opportunity

Advantages of this decision:
* Deleting a previously delivered file was fast by virtue of being able to drop a collection
* Finding all words associated with a file is fast; its just a find on a collection with no filter

Disadvantages:
* Size bloat. Storing documents means additionl text and, currently, object ids. I can prevent the object id issue but not the former.
* No way to take advantage of redundancy between collections. I think there's probably a way to solve that though.

### Project Layout
The layout of this project adheres to the [not-quite-a-standard go project layout](https://github.com/golang-standards/project-layout)

### Design
Ingest and egress are separate modules so that, one day down the line, they could likely be deployed separately in an ALB setup with groups that could scale independently. This problem seemed like it was a good candidate for a load balancer + backend servers since lambdas typically have maximum run time requirements.

The implementation of data storage is abstracted through the Datastore and Document interfaces. As mentioned before, there are a variety of solutions and some day there might be a better one that would be nice to implement without a lot of refactoring.

### Required Packages
#### Logrus
A great logging package that can be customized to do a lot of fun things. I needed logging to see stuff happen.
#### Mux
A great package for improved http handling. Mux provides a more robust way to handle different paths and middleware.
#### Mongo DB Driver
I couldn't possibly write my own :)

### Running The Server
>Note: Its assumed that there is a mongo db server running on localhost that can be reached on port 27017. 
From the project root:
```
CONFIG_FILE_PATH=configs/test.json go run cmd/server.go
```

### Executing Tests
Coverage is really pretty decent. From the root Directory:
```
go test ./... -coverprofile=coverage.out
```
Example Output:
```
% go test ./... -coverprofile=coverage.out
?   	github.com/absolutscottie/bigdocument/cmd	[no test files]
ok  	github.com/absolutscottie/bigdocument/internal/config	0.139s	coverage: 91.7% of statements
?   	github.com/absolutscottie/bigdocument/internal/data	[no test files]
ok  	github.com/absolutscottie/bigdocument/internal/egress	0.386s	coverage: 75.0% of statements
ok  	github.com/absolutscottie/bigdocument/internal/ingest	0.263s	coverage: 66.7% of statements
?   	github.com/absolutscottie/bigdocument/internal/middleware	[no test files]
?   	github.com/absolutscottie/bigdocument/internal/mock	[no test files]
```
