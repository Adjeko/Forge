using Spectre.Console;
using Forge.Rendering;
using Spectre.Console.Rendering;
using Forge.Services.Git;

namespace Forge;

internal static class Program
{
    private static IRenderable? _mainContent;
    private static readonly List<IGitCommand> _commands = new()
    {
        new GitStatusCommand("C:/ADO/CS2"),
        new GitFetchCommand("C:/ADO/CS2")
    };
    private static int _selectedCommandIndex = 0;
    private static bool _isCommandWindowOpen = false;

    private static void Main()
    {
        AnsiConsole.Cursor.Hide();
        Render();

        int lastW = Console.WindowWidth;
        int lastH = Console.WindowHeight;
        while (true)
        {
            if (Console.KeyAvailable)
            {
                var key = Console.ReadKey(true);

                // ESC: entweder Befehlsfenster schließen oder Anwendung beenden
                if (key.Key == ConsoleKey.Escape)
                {
                    if (_isCommandWindowOpen)
                    {
                        _isCommandWindowOpen = false;
                        Render();
                        continue; // zurück in Loop, Anwendung läuft weiter
                    }
                    break; // Anwendung verlassen
                }

                // Ctrl+G öffnet das Befehlsfenster
                if (key.Key == ConsoleKey.G && key.Modifiers.HasFlag(ConsoleModifiers.Control))
                {
                    _isCommandWindowOpen = true;
                    _selectedCommandIndex = Math.Clamp(_selectedCommandIndex, 0, _commands.Count - 1);
                    Render();
                    continue;
                }

                if (_isCommandWindowOpen)
                {
                    switch (key.Key)
                    {
                        case ConsoleKey.UpArrow:
                            _selectedCommandIndex = (_selectedCommandIndex - 1 + _commands.Count) % _commands.Count;
                            Render();
                            continue;
                        case ConsoleKey.DownArrow:
                            _selectedCommandIndex = (_selectedCommandIndex + 1) % _commands.Count;
                            Render();
                            continue;
                        case ConsoleKey.Enter:
                            // ausgewählten Befehl ausführen und Fenster schließen
                            _mainContent = _commands[_selectedCommandIndex].Execute();
                            _isCommandWindowOpen = false;
                            Render();
                            continue;
                    }
                }
            }

            if (Console.WindowWidth != lastW || Console.WindowHeight != lastH)
            {
                lastW = Console.WindowWidth;
                lastH = Console.WindowHeight;
                Render();
            }

            Thread.Sleep(250);
        }
        AnsiConsole.Cursor.Show();
    }

    private static void Render()
    {
        AnsiConsole.Clear();
        var content = _isCommandWindowOpen ? BuildCommandWindow() : _mainContent;
        var layout = MainLayout.Build(startColor: "#FF0000", endColor: "#FFFFFF", mainContent: content);
        AnsiConsole.Write(layout);
    }

    private static IRenderable BuildCommandWindow()
    {
        // Render der Befehlsliste mit Hervorhebung des aktuellen Eintrags
        var table = new Table
        {
            Border = TableBorder.Rounded
        };
        table.AddColumn(new TableColumn("Verfügbare Befehle"));

        for (int i = 0; i < _commands.Count; i++)
        {
            var cmd = _commands[i];
            string label = i == _selectedCommandIndex
                ? $"[yellow]>[/] [bold]{Escape(cmd.Name)}[/]"
                : $"  {Escape(cmd.Name)}";
            table.AddRow(new Markup(label));
        }

        var instructions = new Markup("[grey]Pfeiltasten: Auswahl  |  Enter: Ausführen  |  Esc: Schließen[/]");
        var panel = new Panel(new Rows(table, instructions))
        {
            Header = new PanelHeader("Befehlsfenster"),
            Border = BoxBorder.Rounded
        };
        return panel;
    }

    private static string Escape(string raw) => raw.Replace("[", "[[").Replace("]", "]]" );
}
