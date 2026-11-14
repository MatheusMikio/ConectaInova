# E2E contra a aplicação real

A suíte usa navegador, frontend, API e PostgreSQL reais. Não intercepta nem simula respostas. Além da interface, consulta a API para conferir persistência, relacionamentos e ausência de alterações após requisições rejeitadas.

## Executar

Mantenha o frontend em `http://localhost:5173` e a API em execução, com as migrations aplicadas.

No PowerShell, dentro de `front/`:

```powershell
$env:E2E_ALLOW_WRITES = '1'
bun run test:e2e:real
```

O consentimento explícito é necessário: a suíte cria registros `Teste E2E` que permanecem no banco, porque não existem endpoints de exclusão. Ela não faz reset nem exclui dados existentes. As edições são limitadas aos registros da própria execução.

No Windows, utiliza o Microsoft Edge instalado. Em outro ambiente, instale o navegador com `bunx playwright install chromium`. É possível informar `PLAYWRIGHT_EXECUTABLE_PATH` e `E2E_BASE_URL`.

## Evidências

Cada execução gera `test-results/real-<data>/` com:

- `report.json`: resultado individual, IDs criados e evidências das verificações.
- `trace.zip`: execução do navegador, requisições e snapshots.
- Screenshots dos cenários que falharam e de alguns fluxos.

Para abrir o trace:

```sh
bunx playwright show-trace test-results/real-<data>/trace.zip
```

Os artefatos são ignorados pelo Git. A suíte termina com código diferente de zero quando há cenários falhando ou não executados. Não transforme resultados inesperados em sucessos sem investigar a aplicação.
