<html>
  <head>
    <meta charset="utf-8" />
    <title>🐍 requests</title>
    <script src="https://cdn.jsdelivr.net/npm/@tabler/core@1.0.0-beta11/dist/js/tabler.min.js"></script>
    <link
      rel="stylesheet"
      href="https://cdn.jsdelivr.net/npm/@tabler/core@1.0.0-beta11/dist/css/tabler.min.css"
    />
    <meta name="github:user" content="@sari3l"/>
  </head>
  <body>
    <div class="page-header">
      <div class="container-fluid">
        <div class="row align-items-center">
          <div class="col">
            <div style="font-size: xx-small">
              Reported By:
              <a href="https://github.com/sari3l/requests"
                >https://github.com/sari3l/requests</a
              >
            </div>
            <div class="page-title">Overview</div>
          </div>
          <div class="col-auto ms-auto">{{ .Datetime }}</div>
        </div>
      </div>
    </div>
    <div class="page-body">
      <div class="container-fluid">
        <div class="row" style="margin-bottom: 10px">
          <!-- left column-->
          <div class="col-sm-12 col-md-4 col-lg-5 col-xl-5">
            <div class="row row-cards">
              <div class="col-12">
                <div class="card">
                  <img
                    class="card-img-top"
                    loading="lazy"
                    src="data:image/png;base64,{{ .Snapshot }}"
                  />
                  <div class="card-body">
                    <div class="d-flex">
                      <div>
                        <div>{{ .URL }}</div>
                        <div class="text-muted">{{ .Title }}</div>
                      </div>
                      <a
                        href="{{ .URL }}"
                        class="btn btn-primary ms-auto"
                        >Visit URL
                      </a>
                    </div>
                  </div>
                </div>
              </div>

              <div class="col-12">
                <div class="card">
                  <div class="card-header">
                    <div class="card-title">Console Log</div>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-sm table-vcenter card-table">
                      <thead>
                        <tr>
                          <th>TYPE</th>
                          <th>VALUE</th>
                        </tr>
                      </thead>
                      <tbody style="font-size:12px">
                        {{ range $log := .ConsoleLogs }}
                        <tr>
                          <td class="text-nowrap">{{ $log.Type }}</td>
                          <td class="text-muted">{{ $log.Value }}</td>
                        </tr>
                        {{ end }}
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>

              <div class="col-12">
                <div class="card">
                  <div class="card-header">
                    <div class="card-title">TLS Info</div>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-sm table-vcenter card-table">
                      <thead>
                        <tr>
                          <th class="col-3">SUBJECT CN</th>
                          <th class="col-3">ISSUER CN</th>
                          <th class="col-3">SIG ALGORITHM</th>
                          <th class="col-3">DNS NAMES</th>
                        </tr>
                      </thead>
                      <tbody style="font-size:12px">
                      {{ range $cert := .Certificates}}
                        <tr>
                          <td class="text-nowrap">{{ $cert.Subject.CommonName }}</td>
                          <td class="text-muted">{{ $cert.Issuer.CommonName }}</td>
                          <td class="text-muted">{{ $cert.SignatureAlgorithm }}</td>
                          <td class="text-muted">
                            <ul>
                              {{ range $dnsName := $cert.DNSNames }}
                              <li>{{ $dnsName }}</li>
                              {{ end }}
                            </ul>
                          </td>
                        </tr>
                       {{ end }}
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <!-- right column-->
          <div class="col-sm-12 col-md-8 col-lg-7 col-xl-7">
            <div class="row row-cards">
              <div class="col-12">
                <div class="card">
                  <div class="card-header">
                    <div class="ribbon
                    {{ if and (ge .StatusCode 200) (le .StatusCode 299) }}
                      bg-green
                    {{ else if and (ge .StatusCode 300) (le .StatusCode 399) }}
                      bg-blue
                    {{ else if and (ge .StatusCode 400) (le .StatusCode 500) }}
                      bg-yellow
                    {{ else if and (ge .StatusCode 500) (le .StatusCode 600) }}
                      bg-red
                    {{ end }}
                    ">HTTP {{ .StatusCode }}</div>
                    <div class="card-status-top
                    {{ if and (ge .StatusCode 200) (le .StatusCode 299) }}
                      bg-success
                    {{ else if and (ge .StatusCode 300) (le .StatusCode 399) }}
                      bg-primary
                    {{ else if and (ge .StatusCode 400) (le .StatusCode 500) }}
                      bg-warning
                    {{ else if and (ge .StatusCode 500) (le .StatusCode 600) }}
                      bg-danger
                    {{ end }}
                    "></div>
                    <div class="card-title">Response Headers</div>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-sm table-vcenter card-table">
                      <thead>
                        <tr>
                          <th class="col-sm-4 col-md-3" col-lg-2>KEY</th>
                          <th class="col-sm-8 col-md-9" col-lg-10>VALUE</th>
                        </tr>
                      </thead>
                      <tbody style="font-size:12px">
                        {{ range $key, $value := .Header }}
                          <tr>
                            <td class="text-nowrap">{{ $key }}</td>
                            <td class="text-muted">{{range $v := $value}}{{ $v }}</br>{{end}}</td>
                          </tr>
                        {{ end }}
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>
              <div class="col-12">
                <div class="card">
                  <div class="card-header">
                    <div class="card-title">NetWork Logs</div>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-sm table-vcenter card-table">
                      <thead>
                        <tr>
                          <th class="col-1">TYPE</th>
                          <th class="col-1">CODE</th>
                          <th class="col-2">IP</th>
                          <th class="col-2">ERROR</th>
                          <th class="col-5">URL</th>
                        </tr>
                      </thead>
                      <tbody style="font-size:12px">
                        {{ range $log := .NetworkLogs }}
                        <tr>
                          <td>{{ $log.Type }}</td>
                          <td>
                          {{ if (eq 0 $log.StatusCode) }}
                          {{ else }}
                              {{ if and (ge $log.StatusCode 200) (le $log.StatusCode 299) }}
                                <span class="badge bg-green">{{ .StatusCode }}</span>
                              {{ else if and (ge $log.StatusCode 300) (le $log.StatusCode 399) }}
                                <span class="badge bg-blue">{{ .StatusCode }}</span>
                              {{ else if and (ge $log.StatusCode 400) (le $log.StatusCode 500) }}
                                <span class="badge bg-yellow">{{ .StatusCode }}</span>
                              {{ else if and (ge $log.StatusCode 500) (le $log.StatusCode 600) }}
                                <span class="badge bg-red">{{ $log.StatusCode }}</span>
                              {{ else }}
                                <span class="badge">{{ $log.StatusCode }}</span>
                              {{ end }}
                          {{ end }}
                          </td>
                          <td>{{ $log.IP }}</td>
                          <td>{{ $log.Error }}</td>
                          <td><a href="{{ $log.URL }}" target="_blank">{{ $log.URL }}</a></td>
                        </tr>
                        {{ end }}
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="row" style="margin-bottom: 10px">
          <div class="col-12">
            <div class="card">
              <div class="card-header">
                <div class="card-title">DOM DUMP</div>
              </div>
              <div class="card-body">
                <pre>{{ .DOM }}</pre>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </body>
</html>
