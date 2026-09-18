# MDCD Playground — Java

Maven multi-module project + JUnit 5. Each lecture is a `days/dayXX/`
module: classes with `TODO`s to implement, and JUnit test classes that
check them. Maven downloads JUnit automatically on first run.

## Setup

JDK 17+, Maven, and git.

- **Linux**: `sudo apt install openjdk-17-jdk maven git`
- **macOS**: `brew install openjdk@17 maven` — then either symlink it
  (`sudo ln -sfn "$(brew --prefix openjdk@17)/libexec/openjdk.jdk" /Library/Java/JavaVirtualMachines/openjdk-17.jdk`)
  or add `export JAVA_HOME="$(brew --prefix openjdk@17)/libexec/openjdk.jdk/Contents/Home"`
  to your shell profile
- **Windows**: install a JDK (e.g. [Eclipse Temurin 17](https://adoptium.net/))
  and [Maven](https://maven.apache.org/download.cgi) (or use IntelliJ IDEA /
  Eclipse, both of which bundle Maven support and can open this folder
  directly as a project)

## Build and test

From the `playground/java/` directory, on any OS:

```
mvn test
```

Runs every module. For just one day: `mvn -pl days/day02 test`.

A fresh clone builds cleanly but **fails** most tests — expected, the
`TODO`s aren't filled in yet. Edit `days/dayXX/src/main/java/**/*.java`;
leave `src/test/` alone.
