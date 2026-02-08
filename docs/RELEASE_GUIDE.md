# 🚀 Aviator Release Guide

Questa guida spiega come rilasciare una nuova versione ufficiale di Aviator utilizzando il sistema di automazione basato su **GitHub Actions**.

## 🛠️ Come rilasciare una nuova versione

Il processo è automatizzato: non devi compilare nulla manualmente sul tuo computer. Basta creare un **Tag** su Git.

### 1. Assicurati che tutto sia pushato
Prima di rilasciare, invia tutti i tuoi ultimi cambiamenti sul branch principale:
```bash
git add .
git commit -m "Preparazione release vX.X.X"
git push origin main
```

### 2. Crea un nuovo Tag
Usa il formato `v` seguito dai numeri della versione (es. `v2.1.0`). È importante pushare il tag specificando il remote (solitamente `origin`).
```bash
git tag v2.1.0
git push origin v2.1.0
```

### 3. Cosa succede ora?
1. **GitHub Actions** rileva il nuovo tag pushato su `origin`.
2. Avvia un server Windows virtuale.
3. Installa Go, Node.js e Wails.
4. Compila l'eseguibile in modalità produzione.
5. Crea automaticamente una **Release** nella pagina GitHub del progetto e vi allega l'eseguibile `aviator-wails.exe`.

---

## 🔗 Collegamento Versione (vX.0.0) nell'App

Per far sì che l'app mostri la versione corretta (es. "v2.0.0") invece di un testo fisso, il sistema usa i `ldflags` durante la build.

### Inserimento automatico
Il workflow è configurato per iniettare il nome del tag Git direttamente nel codice Go durante la compilazione:
`wails build -ldflags "-X main.Version=${{ github.ref_name }}"`

### Visualizzazione
- **Desktop**: La versione appare nel log di avvio e può essere richiamata dal frontend Vue tramite il binding `GetAppVersion()`.
- **Web**: La versione è inclusa nella risposta JSON dell'endpoint `/api/info`.

---

## 📝 Note per lo Sviluppo
Se vuoi testare la build localmente con una versione specifica:
```bash
cd aviator-wails
wails build -ldflags "-X main.Version=vTest-Local"
```
