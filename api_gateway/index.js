const app = require("express")();
const proxy = require("http-proxy").createProxyServer();

const PORT = 3000;

app.all("/service-name/*", (req, res) => {
  proxy.web(req, res, { target: "http://localhost:8080" });
});

app.listen(PORT, () => {
  console.log(`gateway active on port ${PORT}`);
});
