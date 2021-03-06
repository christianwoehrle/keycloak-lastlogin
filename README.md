# keycloak-lastlogin



<p align="center">
  <a href="https://github.com/christianwoehrle/keycloak-lastlogin/actions">
    <img src="https://img.shields.io/github/workflow/status/christianwoehrle/keycloak-lastlogin/goreleaser?style=flat-square" alt="Github Actions">
  </a>
<a href="https://godoc.org/github.com/christianwoehrle/keycloak-lastlogin">
    <img src="https://godoc.org/github.com/christianwoehrle/keycloak-lastlogin?status.svg" alt="Github Actions">
  </a>

  <a href="https://goreportcard.com/report/github.com/christianwoehrle/keycloak-lastlogin">
    <img src="https://goreportcard.com/badge/github.com/christianwoehrle/keycloak-lastlogin">
  </a>
  <img src="https://img.shields.io/github/go-mod/go-version/christianwoehrle/keycloak-lastlogin?style=flat-square">
  <a href="https://github.com/christianwoehrle/keycloak-lastlogin/releases">
    <img src="https://img.shields.io/github/release/christianwoehrle/keycloak-lastlogin/all.svg?style=flat-square">
  </a>
</p>

Simple Program to report the last login of users. 



    Usage of lastlogin:

    -dateFrom string
        e.g. 2021-05-10 (default "2021-05-31")

    -password string
        Password to access Keycloak (default "xxx")

    -realm string
        e.g. master (default "master")

    -url string
        Keycloak-URL in the form of https://localhost:8443/ (default "https://localhost:8443/")

    -user string
        Username to access Keycloak (default "admin")



Example output:

/opt/keycloak/bin/keycloak-lastlogin -password $KEYCLOAK_ADMIN_PASSWORD -url http://localhost:8080/ -realm dwpbank -dateFrom "2022-06-30" -log info


swpp1x2t;2022-06-30T13:45:54Z

swpp7x0t;2022-07-01T15:37:48Z

fdbp2x1t;2022-07-01T12:36:28Z
