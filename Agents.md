# Agents

Dieser Text enthält ausschließlich die unumstößlichen, harten Leitplanken.

## Unverrückbare Prinzipien

1. Charm Stack Only: Verwende ausschließlich die offiziellen Charm Libraries: **Bubble Tea**, **Bubbles**, **Bubble Zone**, **Lipgloss**, **Huh**.
2. CLI vor Client-SDK: Externe Werkzeuge/Services grundsätzlich über ihre offiziellen CLI-Programme aufrufen (z.B. `git`, `gh`, `docker`) – keine ersetzenden Client-Bibliotheken, wenn ein stabiler CLI-Aufruf möglich ist.
3. Minimale Dependencies: So wenig Fremdbibliotheken wie möglich. Bevorzugt: Go Standardbibliothek + Charm Stack. Neue Abhängigkeit nur bei klarer, dokumentierter Notwendigkeit (Security / Performance / fehlende Funktion in Stdlib).
4. Transparente Commands: Jeder ausgeführte CLI-Befehl soll für den Nutzer nachvollziehbar/logbar sein (Option zur Anzeige oder Debug-Modus). Keine Intransparenz durch Wrapper, die Semantik verändern.
5. Fehlerkultur: Keine Panics im regulären Ausführungspfad. Fehler werden erfasst, klar klassifiziert und nutzerfreundlich dargestellt. Fail fast bei irreparablen Zuständen – aber kontrolliert.
6. UI Reinheit: Präsentationslogik bleibt in Bubble Tea / Lipgloss Komponenten. Keine Vermischung von Anwendungslogik und Rendering.
7. Ressourcen-Schonung: Prozesse & Subprozesse nur starten, wenn nötig; sauber beenden. Kein Polling ohne konfiguriertes Intervall.
8. Sicherheit vor Komfort: Keine automatische Ausführung destruktiver Befehle (z.B. `git push`, `rm -rf`) ohne ausdrückliche Bestätigung.
9. Konfigurationsobergrenze: Nur wirklich notwendige Konfigurationsoptionen. Default-Werte müssen sinnvoll und sicher sein.
10. Erweiterbarkeit ohne Bruch: Neue Agents dürfen bestehende unumstößliche Prinzipien nicht verletzen.
11. Tests nur auf Anforderung: Es werden keine neuen Tests geschrieben, außer der Nutzer fordert sie ausdrücklich an oder ein kritischer Fehler verlangt Regression-Schutz.
12. Einzelentwickler-Fokus: Priorität hat klare Verständlichkeit des Codes, einfache Erweiterbarkeit und geringe kognitive Last gegenüber komplexer Architektur oder überoptimierter Abstraktion.
13. Deutsche Kommentare: Quellcode-Kommentare werden ausschließlich in deutscher Sprache verfasst.
14. Deutsche Ausgabe: Der Agent verwendet in seiner Benutzer-Ausgabe ausschließlich die deutsche Sprache (UI, Logs, Fehlerhinweise), außer ein expliziter Nutzerwunsch verlangt Mehrsprachigkeit.
15. Model-Struktur: Jedes UI-Model besteht zwingend aus einem `struct`, der alle dafür notwendigen Zustands-Variablen enthält, sowie genau den Methoden `Init`, `Update` und `View`. Neue, eigenständige Anzeige- oder Interaktions-Elemente werden so klein wie möglich gehalten und in eigene, separierte Model-Implementierungen ausgelagert (Kompositionsansatz statt monolithischer Wachstum). Keine überflüssigen Methoden oder Felder.

## Nicht-Ziele (implizit ausgeschlossen)
- Nutzung alternativer UI-Frameworks außerhalb Charm Stack.
- Heavy ORM / komplexe Abstraktionslayer ohne zwingenden Grund.
- Undokumentierte Hintergrundjobs.
- Verdeckte Telemetrie oder Tracking.

## Änderung dieser Datei
Änderungen nur, wenn ein Prinzip nachweislich verletzt wird oder eine neue, zwingende Grundlage entsteht. Diskussion erforderlich; kein Einzel-Commit ohne Review.
