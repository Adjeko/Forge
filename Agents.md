# AGENTS.md

## Zweck
Dieses Repository enthält eine reine Terminal User Interface (TUI) Anwendung. Ziel ist es, eine robuste, wartbare und einheitliche Konsolenoberfläche bereitzustellen.

## Verwendete Technologie
- Primärer (und einziger) UI-Stack: **Spectre.Console** / **Spectre.Console.Cli** (im Folgenden kurz: *Spectre.NET*)
- Es dürfen **keine weiteren UI-/TUI-Bibliotheken** eingebunden werden (kein curses, keine eigenen ANSI-Helfer, kein Figlet außerhalb von Spectre).

## Grundprinzipien
1. Verwende Spectre.NET konsequent für sämtliche Darstellung (Layouts, Tables, Trees, Progress, Panels, Markup, Coloring, Prompts, Command Handling).
2. Implementiere **nichts manuell**, wenn Spectre.NET bereits eine vergleichbare oder ausreichende Funktionalität bietet.
3. Bevor eigener Code für Formatierung, Rendering, Input-Flow oder Command-Parsing entsteht: Dokumentiere kurz, warum die Spectre-Funktion nicht reicht. (Falls später nötig, Issue anlegen.)
4. Bevorzugt deklarative/kompositorische Nutzung der vorhandenen Spectre-Objekte statt "Low-Level"-String-Bastelei.
5. Keine direkten ANSI-Escape-Sequenzen schreiben (kein `\u001b[...]`). Farb- und Stilsteuerung ausschließlich über Spectre-Markup oder API.
6. Output soll deterministisch, strukturiert und klar gegliedert sein. Wiederverwendbare Render-Komponenten kapseln.
7. Commands folgen den Spectre.Console.Cli-Konventionen (Settings-Klassen, Validierung über Attributes / Overrides / TypeConverter, konsistente Hilfeausgabe).

## Do / Don't
**Do:**
- `AnsiConsole.Markup`, `AnsiConsole.Write`, `Table`, `Tree`, `Progress`, `Panel`, `Rule`, `BarChart`, `FigletText`, `Status`, `Prompt` verwenden.
- Command-Struktur über `CommandApp` und `ICommand`/`Command<TSettings>` organisieren.
- Wiederverwendbare Komponenten als kleine Hilfsklassen / Renderer kapseln (sofern sie ausschließlich Spectre-Aufrufe orchestrieren).

**Don't:**
- Eigene Parsing- oder Rendering-Layer bauen, wenn Spectre.Console.Cli / Rendering API genügt.
- Externe Libraries für CLI Styling / Parsing hinzufügen (z.B. System.CommandLine, Colorful.Console, Sharprompt, etc.).
- Rohes Schreiben auf `Console.*` statt über `AnsiConsole` (Ausnahme: sehr früher Bootstrap vor Spectre-Init – nach Möglichkeit vermeiden).
- Manuelle Threading-Loops für Spinner/Progress – Spectre Progress/Status benutzen.

## Erweiterungen / Zukunft
Falls Anforderungen auftreten, die Spectre.NET nicht abdeckt:
1. Prüfen: Existiert ein offizielles Feature oder Issue im Spectre.Console Repo?
2. Issue im eigenen Repo mit Label `needs-spectre-extension` erstellen.
3. Nur wenn absolut nötig: Dünne, klar begrenzte Erweiterungsschicht (z.B. Adapter) erstellen – ohne Spectre-Kern zu umgehen.

## Qualitätsrichtlinien
- Keine doppelten Farbkodierungen; zentrale Farb-/Style-Definitionen (später ggf. Theme-Datei).
- Testbarkeit: Logik (Domain, Parsing, Validierung) von Render-Code trennen.
- Konsistente Konsolenbreite berücksichtigen (Fallback gracefully bei kleiner Breite, keine Annahmen über 120+ Spalten fest einbacken).

## Namenskonventionen (Vorschlag)
- Ordner `Commands/`, `Rendering/`, `Services/`, `Models/`.
- Render-Helfer enden auf `Renderer` oder `View` (z.B. `TaskSummaryRenderer`).
- Commands enden auf `Command` (z.B. `BuildCommand`).

## Entscheidungslog / Architektur
Wenn von obigen Regeln abgewichen wird: Kurz im README oder in `docs/ARCHITECTURE.md` festhalten (Rationale + Datum + Autor).

## Kurzfassung für neue Mitwirkende
Diese App ist eine TUI. Benutze ausschließlich Spectre.Console (Spectre.NET). Baue keine Features nach, die es dort schon gibt. Keine Zusatz-UI-Libraries. Kein manueller ANSI-Output.

---
Letzte Aktualisierung: 2025-11-06
