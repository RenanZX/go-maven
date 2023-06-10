# GO Maven POM Generator

Is a sample pom generator, that search jar libs in maven repository and generate a ```pom.xml``` file from it.

# Run the project

Follow the steps:

1. You need to put a ```lib``` folder with your ```.jar``` dependencies at the root of the project.
2. Change ```param.config``` file to your desired settings:
   1. ```modelVersion``` is the model version of your project
   2. ```groupId``` is the groupId of your project
   3. ```artifactId``` the artifactId of the project
   4. ```version``` version of the project
   5. ```name``` is the name of the project
   6. ```includeLibsNotFound``` it's going to include libs in your project that is not encountered at maven repository from jar folder.
3. Type and run ```./go-maven``` or ```go run main.go``` in your terminal.
4. It will be generated a ```pom.xml``` in the root folder.