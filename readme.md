

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