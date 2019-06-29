Stagevariable

on http handler if you not in environement lambda all road has prefixed by ```/:__stage__```
when request oncomming the stage variable is stored in event.RequestContext.Stage 


if you need to pass a stageVariables to lambda with http handler use the header ```Stagevariable_{var_name}```
exemple: 
```
request header:
Stagevariable_foo=bar

lambda handler:
print(req.StageVariables["foo"]) // output: bar 

``` 