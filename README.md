# gradle-resolver-redirector
Run this app on your machine and add following piece in repository section of build.gradle file :
```
maven{
    url 'http://your-machine-ip:10010'
}
```

It will listen to 10010 port by default.

This app will check following repositories for libraries:

> https://dl.google.com/dl/android/maven2

> https://jcenter.bintray.com/

> https://repo1.maven.org/maven2

