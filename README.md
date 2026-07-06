# Go Scraper Learning

Este proyecto está pensado para aprender a hacer web scrapping en Go usando la librería `goquery`. Para probar se usará la página web [Books To Scrape](https://books.toscrape.com/), la cual está pensada expresamente para realizar scrapping sobre ella.

Se implementarán sistemas adicionales como por ejemplo sistema de logging en archivos locales `.log`, guardado de los datos recogidos en bases de datos e incluso logging pero a través de utilidades como los bots de telegram.

NOTA: este software ha sido desarrollado en Linux, es posible que algún comando no funcione en otros SO. (Lo dudo, pero por si acaso :D)

## ¿Cómo ejecutar el proyecto?

Si sólo se quiere ejecutar, la solución es el siguiente comando estando en la raíz del proyecto:

```shell
make run
```

ó

```go
go run cmd/scraper/main.go
```

Internamente, el comando `make run` ejecuta ese segundo comando (aunque inyectando variables).

---

Si lo que se quiere hacer es compilar a binario (construir), se podrá ejecutar lo siguiente:

```shell
make build
```

ó

```go
go build -o build/scraper cmd/scraper/main.go
```

Al igual que con la ejecución normal del programa, el comando `make build` ejecuta internamente ejecuta ese mismo comando, pero inyectando variables.

---

Una vez el binario esté compilador, se podrá ejecutar con el siguiente comando:

```shell
make run-binary
```

ó

```shell
./build/scraper
```

## ¿Cómo configurar tu propia base de datos y bot de telegram?
