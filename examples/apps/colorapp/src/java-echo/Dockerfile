FROM docker.io/openjdk:21
EXPOSE 2701
WORKDIR /
ADD target/java-echo-0.0.1-SNAPSHOT.jar app.jar
#启动应用固定命令
ENTRYPOINT ["java","-jar"]
#执行jar包命令
CMD ["app.jar"]





