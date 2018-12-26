<!DOCTYPE html>
<html lang="en">

<head>
  <title>{{.PageTitle}}</title>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css">
  <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.3/umd/popper.min.js"></script>
  <script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.1.3/js/bootstrap.min.js"></script>
  <link href="https://fonts.googleapis.com/css?family=Didact+Gothic|Ubuntu" rel="stylesheet">
  <link rel="shortcut icon" type="image/png" href="/static/batman.png" />

</head>
<style>
  body {
    font-family: 'Didact Gothic', sans-serif;
    max-width: 1000px;
    margin: auto;
  }

  .debuginfo {
    font-family: 'Ubuntu', sans-serif;
  }
  
  .notimportant{
    color: white;
  }

</style>

<body>

  <div class="jumbotron text-info">
    <h2 class="notimportant">All important things in life start with a white Screen</h2>
    <h4 class="text-right"> -- <small>{{.AppName}}:{{.ApiVersion}} / Host:{{.Host}} </small>  </h4>
    <p>Serving from Host : {{.Host}}</p>
    {{ if .message }}
      <p class="alert alert-info">{{.message}}</p>
    {{end}}
  </div>

  <div class="card font-small">
    <div class="card-header">
      <h2>Request headers</h2>
    </div>
    <div class="card-body debuginfo">
      <table class="table">
        <thead>
          <tr>
            <th scope="col">Name</th>
            <th scope="col">Value</th>
          </tr>
        </thead>
        <tbody>
          {{ range $key, $value := .Request.Headers }}
          <tr>
            <td>{{$key}}</td>
            <td>
              <span class="d-inline-block text-truncate" style="max-width: 700px;">
                {{$value}}
              </span>
            </td>
          </tr>
          {{ end }}
        </tbody>
      </table>
    </div>
  </div>
  {{if .Environment}}
  <div class="card font-small">
    <div class="card-header">
      <h2>Environment Vars</h2>
    </div>
    <div class="card-body debuginfo">
      <table class="table">
        <thead>
          <tr>
            <th scope="col">Name</th>
            <th scope="col">Value</th>
          </tr>
        </thead>
        <tbody>
          {{ range $key, $value := .Environment }}
          <tr>
            <td>{{$key}}</td>
            <td>
              <span class="d-inline-block text-truncate" style="max-width: 700px;">
                {{$value}}
              </span>
            </td>
          </tr>
          {{ end }}
        </tbody>
      </table>
    </div>
  </div>
  {{ end }}
  <footer class="page-footer font-small blue pt-4">
    <div class="footer-copyright text-center py-3">GitHub:
      <a href="https://github.com/ajarv/go-app-docker.git"> https://github.com/ajarv/go-app-docker.git</a>
    </div>
  </footer>
</body>

</html>