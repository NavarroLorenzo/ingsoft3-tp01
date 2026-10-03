const fs = require("node:fs");

const reportPath = process.argv[2] || "coverage/report/coverage-summary.json";
const { total } = JSON.parse(fs.readFileSync(reportPath, "utf8"));

if (!total?.lines?.total) {
  throw new Error("El reporte de frontend no midió archivos. Revisar coverage.include.");
}

console.log([
  "## Cobertura del frontend",
  "",
  "| Métrica | Resultado |",
  "|---|---:|",
  `| Líneas | ${total.lines.pct}% |`,
  `| Ramas | ${total.branches.pct}% |`,
  `| Funciones | ${total.functions.pct}% |`,
  `| Sentencias | ${total.statements.pct}% |`
].join("\n"));
