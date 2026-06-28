package service

import _ "net/http" // want `service package must not import transport framework "net/http"\.`

type BadManager struct{} // want `type "BadManager" uses forbidden suffix\.`
