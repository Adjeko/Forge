

# Usecases
❓muss noch geklärt werden
❕Geklärt, muss implementiert werden
✅ fertig implementiert



## ❕main in den aktuellen Branch holen

| Schritt | Ablauf | Input | Notizen |
| :--- | :---: | :---: | ---: |
| 1 | aktuellen Branchnamen speichern | RepoPfad | |
| 2 | git checkout main | RepoPfad | |
| 3 | git fetch | RepoPfad | |
| 4 | git pull | RepoPfad | (main) |
| 5 | git checkout [Branchname] | RepoPfad, Branchname | |
| 6 | git merge main | RepoPfad | |
| 7 | Projekt bauen | RepoPfad | optional |

## ❕alles clean bauen

| Schritt | Ablauf | Input |Notizen |
| :--- | :---: | :---: | ---: |
| 1 | Movisuite Client schließen | Prozessname | |
| 2 | Movisuite Server schließen | Prozessname | |
| 3 | Visual Studio schließen | Prozessname | anderen IDEs ? |
| 4 | Projekt bauen | RepoPfad | |
| 5 | Visual Studio öffnen | SolutionPfad | |

## ❕Moduldaten aktualisieren

| Schritt | Ablauf | Input |Notizen |
| :--- | :---: |:---: | ---: | 
| 1 | update Moduldaten | ModuldatenPfad | |
| 2 | link Moduldaten | RepoPfad, ModuldatenPfad | admin! |
| 3 | Projekt bauen | RepoPfad | |

## ✅ Start der TUI (SEW Forge)

| Schritt | Ablauf | Kommando | Hinweis |
| :--- | :---: | :---: | ---: |
| 1 | Abhängigkeiten prüfen | `go mod tidy` | Lädt fehlende Charm Bibliotheken |
| 2 | Anwendung starten | `go run .` | Zeigt "SEW Forge" an |
| 3 | Beenden | `Ctrl+C` oder `q` | Sauberer Programmabbruch |

### Beschreibung
Die minimalistische TUI nutzt Bubble Tea + Lipgloss und zeigt derzeit nur den Titel "SEW Forge" in hervorgehobener Formatierung. Sie dient als Ausgangspunkt für weitere Funktionen.

### Nächste mögliche Erweiterungen (nur bei Bedarf)
- Lade-Animation beim Start.
- Konfigurierbare Farben über eine TOML-Datei.
- Integration weiterer Bubbles (Listen, Textinput) für Benutzerinteraktion.

Alle zukünftigen UI-Ausgaben und Kommentare bleiben gemäß Leitplanken auf Deutsch.