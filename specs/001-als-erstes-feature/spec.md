# Feature Specification: Forge TUI Task Runner Core Interface & Workflow Execution (Initial Surface)

**Feature Branch**: `001-als-erstes-feature`  
**Created**: 2025-10-10  
**Status**: Draft  
**Input**: User description (German original): "Als erstes Feature sollen große Teile der Oberfläche entwickelt werden..."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Einzelnes primitives Kommando ausführen (Priority: P1)

Ein Entwickler öffnet das Programm und wählt ein einzelnes vordefiniertes (hartcodiertes) primitives Kommando (z.B. "git status") über die Oberfläche aus und führt es aus. Die Ausgabe erscheint im linken Terminalbereich vollständig sichtbar, farblich eindeutig dem Schritt zugeordnet.

**Why this priority**: Fundamentale Basis – ohne Einzelausführung kein Nutzen.

**Independent Test**: Programm starten → "git status" auswählen → Ausführen → Farbausgabe sichtbar, Exit-Ergebnis ersichtlich, kein Crash.

**Acceptance Scenarios**:

1. **Given** Programm gestartet, **When** Nutzer wählt "git status" und startet, **Then** Ausgabe erscheint farblich markiert; Fehlertext rot falls Fehler.
2. **Given** Kommando liefert keine Ausgabe, **When** ausgeführt, **Then** leerer gekennzeichneter Abschnitt mit Erfolgsstatus.

---

### User Story 2 - Sequenziellen Ablauf aus primitiven Kommandos ausführen (Priority: P2)

Mehrere primitive Kommandos werden zu einem Ablauf (Workflow) verkettet und einmalig nacheinander ausgeführt; Fehler stoppt Ablauf (Fail-Fast).

**Why this priority**: Haupt-Produktivitätsgewinn durch Automatisierung wiederkehrender Sequenzen.

**Independent Test**: Ablauf (status → fetch → log -1) definieren und starten; Reihenfolge, Farbcodierung, Abbruch bei Fehler prüfen.

**Acceptance Scenarios**:

1. **Given** Ablauf mit 3 Schritten, **When** gestartet, **Then** Ausgaben erscheinen in definierter Reihenfolge mit unterschiedlichen Farben.
2. **Given** Schritt 2 erzeugt Fehler (Exit Code ≠ 0), **When** Ablauf läuft, **Then** Schritt 2 rot, Ablauf stoppt, Schritt 3 nicht ausgeführt.

---

### User Story 3 - Status-LEDs & Fortschrittsleiste (Priority: P3)

Während der Ausführung visualisiert rechter Bereich Schrittzustände mittels LEDs; Footer zeigt zweizeiligen Fortschrittsbalken (Name + Gradient).

**Why this priority**: Erhöht Transparenz & reduziert kognitive Last.

**Independent Test**: 4-Schritt-Ablauf starten, während Lauf LED-Zustände & Fortschritt beobachten.

**Acceptance Scenarios**:

1. **Given** Ablauf läuft, **When** Schritt 2 aktiv, **Then** LED Schritt 2 aktiv, vorherige erfolgreich, folgende neutral.
2. **Given** Fehler tritt auf, **When** Schritt fehlschlägt, **Then** LED des Schritts rot, Fortschrittsbalken zeigt Endstatus (Failure) an.

---

### User Story 4 - Automatische Farbzuteilung (Priority: P4)

Automatische, konsistente Farbzuteilung je Schritt; Fehlerausgaben immer rot.

**Why this priority**: Schnellere visuelle Differenzierung bei mehreren Schritten.

**Independent Test**: 5-Schritt-Ablauf ausführen, unterschiedliche Farben prüfen; Fehler überschreibt Farbe.

**Acceptance Scenarios**:

1. **Given** 5 Schritte, **When** ausgeführt, **Then** jeder Schritt hat eindeutige Farbe (Palette erlaubt).
2. **Given** Schritt 3 Fehler, **When** ausgeführt, **Then** gesamte Ausgabe Schritt 3 rot.

---

### Edge Cases

- Leeres Kommando → Validierungsfehler, keine Ausführung.
- Extrem viel Output → unlimitierter Puffer; Performance-Test (50k Zeilen) muss ohne spürbaren Freeze (>500 ms Input-Reaktionszeit) bestehen.
- Ablauf mit 0 Schritten → Start blockiert mit Hinweis.
- Schritt ohne neue Ausgabe > 60s → Status „hängend“ angezeigt.
- Farbpalette erschöpft → Wiederverwendung erlaubt; Schrittindex + LED sichern Unterscheidbarkeit.
- Fensterbreite < 100 Zeichen → Degradiertes Layout (minimaler LED-Bereich oder Toggle), Fokus auf Output.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST render eine Vollbild-TUI mit Header (4 Zeilen: obere Info + 3-Zeilen ASCII Forge-Schriftzug) und Footer (Progress-Bereich) sowie Mittelteil zweigeteilt (links 2/3 Output, rechts 1/3 Status-LEDs).
- **FR-002**: System MUST display den Farbverlauf (#B50013 → #CCCCCC) für das Wort "Forge" im Header und für aktive Fortschrittsbalken.
- **FR-003**: System MUST allow Auswahl & Ausführung eines einzelnen primitiven Kommandos (hartcodierte Liste initial enthält mindestens "git status", "git fetch", "git log -1").
- **FR-004**: System MUST execute primitive Kommandos sequenziell innerhalb eines Ablaufs und nacheinander deren Ausgaben anzeigen.
- **FR-005**: System MUST color-code jede Schritt-Ausgabe mit einer zugewiesenen Farbe; Fehlerschritte überschreiben mit Rot.
- **FR-006**: System MUST stop remaining steps of a workflow after first failing step unless user configured continuation (Konfiguration zukünftiges Feature; aktuell: immer stoppen).
- **FR-007**: System MUST show Status-LEDs rechts mit Zuständen: Pending, Running, Success, Failed für jeden definierten Schritt eines gestarteten Ablaufs.
- **FR-008**: System MUST show im Footer einen zweizeiligen Fortschrittsbereich: (Zeile 1) Ablaufname, (Zeile 2) fortschreitender Balken (Gradient gefüllt proportional erledigten Schritten).
- **FR-009**: System MUST validate before start: Kein Ablaufstart wenn 0 Schritte; Fehlermeldung sichtbar ohne Absturz.
- **FR-010**: System MUST handle large output streaming without freezing UI (pufferisierte Anzeige, mindestens letzte 5000 Zeilen sichtbar). [Assumption]
- **FR-011**: System MUST indicate completion state (Success / Failed) im Footer nach Ende.
- **FR-012**: System MUST assign deterministic colors to steps from predefined palette (z.B. Gelb, Blau, Grün, Magenta, Cyan, Weiß… Rot reserviert für Fehler).
- **FR-013**: System SHOULD degrade gracefully on narrow window widths (unter Mindestbreite: Status-LED Bereich reduziert oder umschaltbar). [NEEDS CLARIFICATION: Schwellenwert für Layout-Umschaltung]
- **FR-014**: System MUST ensure that header and footer remain fixed while middle content scrolls.
- **FR-015**: System MUST record exit status per Schritt für LED-Zustand.
- **FR-016**: System MUST prevent execution of empty or whitespace-only commands (Validierungsfehler anzeigen).
- **FR-017**: System MUST indicate hanging step status after 60 seconds of no output.

### Key Entities

- **Primitive Command**: Ein vordefiniertes ausführbares CLI-Kommando (Name, Anzeige-Label, zugewiesene Farbe, letzter Exit-Code).
- **Workflow (Ablauf)**: Sequenz von Primitive Commands (Name, Schritt-Liste, Status: Pending/Running/Success/Failed, Start-/Endzeit, Fehler-Schritt-ID falls abgebrochen).
- **Execution Step**: Laufzeitinstanz eines Primitive Command innerhalb eines Workflows (Referenz auf Kommando, Start-/Endzeit, Exit-Code, Output-Puffer, Farbe, Status).
- **Status LED**: Visuelle Abstraktion eines Execution Step Zustands (Symbol + Farbe pro Zustand).
- **Progress Bar**: Darstellung des prozentualen Fortschritts (# erledigte Schritte / Gesamt). Enthält Titel (Workflow-Name) und Gradient-Füllung.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Einzelnes Kommando in ≤5s nach Programmstart ausführbar inkl. sichtbarer Ausgabe.
- **SC-002**: 5-Schritt-Ablauf zeigt farbige Ausgaben ohne UI-Freeze; Eingabereaktion ≤200 ms während Streaming.
- **SC-003**: Fehler stoppt Ablauf; Markierung aller verbleibenden Schritte innerhalb ≤1s.
- **SC-004**: Fortschrittsanzeige aktualisiert ≤250 ms nach Schrittabschluss; erreicht 100% bei Ende.
- **SC-005**: 100% Testpersonen erkennen Status jedes Schritts korrekt (Testlauf mit 5 Personen).
- **SC-006**: Layout ≥100 Zeichen ohne Überlappung; <100 Zeichen: definierte Degradierung ohne Verlust kritischer Statusinfos.

### Assumptions

- Unlimitierter Output-Puffer (Bewusste Entscheidung; Performance überwachen).
- Timeout hängender Schritt: 60s fix im MVP.
- Mindestbreite für Zweispalten: 100 Zeichen.
- Einzelner Entwickler als Benutzerrolle.
- Fail-Fast vereinfacht Fehlerdiagnose.
- Feste Farbpalette im MVP (keine Konfiguration).

## Clarifications

### Session 2025-10-10

- Q: Maximale Puffertiefe für Ausgabe? → A: Unbegrenzt (kein Limit)
- Q: Timeout-Länge für hängende Schritte? → A: 60 Sekunden
- Q: Mindestbreite für zweispaltiges Layout? → A: 100 Zeichen
