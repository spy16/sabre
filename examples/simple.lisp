; define a value binding
(def π 3.1412)

; define a lambda function
(def echo (fn* [msg] msg))

; index access for a vector
([1 2 3] 0)

(def hello (fn*
    ([] "Hello")
    ([name] (let* [user-name name
                  msg (str "Hello " user-name "!")]
                msg))
))

(def defn (fn* [name args body]
    (do
        (if (not (symbol? name))
            (throw "name must be symbol, not " (type name)))
        (if (not (vector? args))
            (throw "args must be a vector, not " (type args)))
        (eval `(def ~name (fn ~args ~body))))))
