<!DOCTYPE html>
<html lang="en">

<head>
  <!-- Required meta tags -->
  <meta charset="utf-8">
  <meta content="initial-scale=1, shrink-to-fit=no, width=device-width" name="viewport">

  <!-- CSS -->
  <!-- Add Material font (Roboto) and Material icon as needed -->
  <link href="https://fonts.googleapis.com/css?family=Roboto:300,300i,400,400i,500,500i,700,700i|Roboto+Mono:300,400,700|Roboto+Slab:300,400,700"
    rel="stylesheet">
  <link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">

  <!-- Add Material CSS, replace Bootstrap CSS -->
  <link href="static/css/material.min.css" rel="stylesheet">
  <style>
    body {
    font-family: 'Didact Gothic', sans-serif;
    max-width: 1200px;
    margin: auto;
  }

  .debuginfo {
    font-family: 'Ubuntu', sans-serif;
  }

  .notimportant {
    color: white;
    background-color: aliceblue;
  }

  .sayings {
    font-size: x-small;
  }
</style>
</head>

<body>
  <nav class="navbar navbar-expand-md navbar-dark fixed-top bg-dark">
    <a class="navbar-brand" href="#">Leego Parts</a>
    <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarsExampleDefault"
      aria-controls="navbarsExampleDefault" aria-expanded="false" aria-label="Toggle navigation">
      <span class="navbar-toggler-icon"></span>
    </button>

    <div class="collapse navbar-collapse" id="navbarsExampleDefault">
      <ul class="navbar-nav mr-auto">
        <li class="nav-item active">
          <a class="nav-link" href="#">Home <span class="sr-only">(current)</span></a>
        </li>
        <li class="nav-item">
          <a class="nav-link" href="#">Link</a>
        </li>
        <li class="nav-item">
          <a class="nav-link disabled" href="#">Disabled</a>
        </li>
        <li class="nav-item dropdown">
          <a class="nav-link dropdown-toggle" href="http://example.com" id="dropdown01" data-toggle="dropdown"
            aria-haspopup="true" aria-expanded="false">Dropdown</a>
          <div class="dropdown-menu" aria-labelledby="dropdown01">
            <a class="dropdown-item" href="#">Action</a>
            <a class="dropdown-item" href="#">Another action</a>
            <a class="dropdown-item" href="#">Something else here</a>
          </div>
        </li>
      </ul>
      <form class="form-inline my-2 my-lg-0">
        <input class="form-control mr-sm-2" type="text" placeholder="Search" aria-label="Search">
        <button class="btn btn-outline-success my-2 my-sm-0" type="submit">Search</button>
      </form>
    </div>
  </nav>
  <br />
  <br />
  <br />
  <br />
  <main role="main">

    <div class="container jumbotron text-info">
      <h2 class="notimportant">All important things in life start with white</h2>
      <h4 class="text-right text-danger font-italic font-weight-light">
        <span>--&gt;</span>
        <span class="badge badge-secondary">
          {{.AppName}} | {{.ApiVersion}}
        </span>
        <span class="badge badge-primary">
          Host {{.Host}}
        </span>
        {{ if .Environment.CF_INSTANCE_INDEX }}
        <span class="badge badge-info">
          CF Instance {{.Environment.CF_INSTANCE_INDEX}}
        </span>
        {{end}}
      </h4>
      {{ if .message }}
      <p class="alert alert-info">{{.message}}</p>
      {{end}}
      {{ if .warning }}
      <p class="alert alert-warning">{{.warning}}</p>
      {{end}}
      <div class="float-right">
        <a href="/" class="btn btn-dark" role="button">Start over</a>
        <a href="/redis" class="btn btn-success" role="button">Hit redis</a>
        <a href="/die" class="btn btn-danger" role="button">Kill server</a>
      </div>
    </div>
    <br />
    <br />
    <br />
    <div class="container card font-small sayings notimportant">
      <div class="card-header">
        <h2>...</h2>
      </div>
      <div class="row">
        <div class="col-sm-4">
          <h4>Remember there are two type of addictions</h4>
          <p></p>
          <p> The first where euphoria comes first followed by a LOT of questioning.<br />
            .. friday night , party .. you drink yourself to brim and Saturday morning you wake up in hangover and pain
            questioning how did you end up like this.
          </p>
          <p>The second type where questioning comes first and euphoria follows.<br />
            .. 5am .. you are outside running in the cold, questioning how did you end up on the street like this. At 7
            am
            you are back in your home , drinking coffee .. pure euphoria.
          </p>
          <br />
          <h3>Choose your addictions wisely!</h3>
        </div>
        <div class="col-sm-8">
          <div>
            <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" id="Capa_1"
              x="0px" y="0px" width="256px" height="256px" viewBox="0 0 29.904 29.904" style="enable-background:new 0 0 29.904 29.904;"
              xml:space="preserve">
              <g>
                <g>
                  <path d="M28.405,14.7c-0.479,0-0.897,0.228-1.172,0.576c-1.56-1.127-4.992-2.994-7.975-0.271c0,0-3.021-4.168-0.982-7.569    c0.246,0.178,0.547,0.286,0.875,0.286c0.827,0,1.5-0.671,1.5-1.5s-0.673-1.5-1.5-1.5c-0.828,0-1.502,0.671-1.502,1.5    c0,0.168,0.032,0.327,0.084,0.478c-2.141,0.819-5.836,2.858-6.39,7.307c0,0-3.429-4.541-8.573-1.594    c-0.265-0.425-0.732-0.711-1.27-0.711c-0.829,0-1.501,0.672-1.501,1.5s0.672,1.5,1.501,1.5c0.828,0,1.499-0.672,1.499-1.5    c0-0.047-0.01-0.091-0.014-0.137c1.794,0.14,4.67,1.726,5.461,10.151l0.09,0.688c0,0.707,2.858,1.279,6.382,1.279    c3.526,0,6.383-0.574,6.383-1.279c0,0,0.229-5.78,5.611-7.623c0.041,0.791,0.688,1.423,1.491,1.423c0.83,0,1.5-0.673,1.5-1.5    C29.907,15.371,29.235,14.7,28.405,14.7z"
                    fill="white" />
                </g>
              </g>
            </svg>
          </div>
          <br />
        </div>
      </div>
    </div>
    <br />

    <div class="container card font-small">
      <div class="card-header">
        <h2>Request headers</h2>
      </div>
      <div class="card-body debuginfo">
        <table class="table table-sm">
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

    <br />
    {{if .Environment}}
    <div class="container card font-small">
      <div class="card-header">
        <h2>Environment Vars</h2>
      </div>
      <div class="card-body debuginfo">
        <table class="table table-sm">
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
                <span class="d-inline-block" style="max-width: 700px;">
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
  </main>

  <footer class="page-footer font-small blue pt-4">
    <div class="footer-copyright text-center py-3">GitHub:
      <a href="https://github.com/ajarv/go-app-docker.git"> https://github.com/ajarv/go-app.git</a>
    </div>
  </footer>

  <!-- Optional JavaScript -->
  <!-- jQuery first, then Popper.js, then Bootstrap JS -->
  <script src="https://code.jquery.com/jquery-3.3.1.slim.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.3/umd/popper.min.js"></script>
  <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.1.1/js/bootstrap.min.js"></script>

  <!-- Then Material JavaScript on top of Bootstrap JavaScript -->
  <script src="static/js/material.min.js"></script>

</body>

</html>