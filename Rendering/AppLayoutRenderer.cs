using Spectre.Console;
using Spectre.Console.Rendering;

namespace Forge.Rendering;

/// <summary>
/// Erzeugt das Grundlayout der Anwendung: 5 Zeilen Header, flexibler Inhalt, 2 Zeilen Footer.
/// Kapselt ausschließlich Spectre.Console Aufrufe gemäß Richtlinien.
/// </summary>
internal static class AppLayoutRenderer
{
    /// <summary>
    /// Baut ein neues Layout für die aktuelle Render-Periode.
    /// </summary>
    public static Layout Build()
    {
        var layout = new Layout("root")
            .SplitRows(
                new Layout("header").Size(5), // Header jetzt fünfzeilig für ASCII-FORGE
                new Layout("content"),
                new Layout("footer").Size(2)
            );

        layout["header"].Update(BuildHeader());
        layout["content"].Update(BuildContent());
        layout["footer"].Update(BuildFooter());

        return layout;
    }

    private static IRenderable BuildHeader()
    {
        // Dreizeilige Darstellung des Wortes FORGE aus den Zeichen █, ▀, ▄
        // Konvention: Zeile 1 (Top) vorwiegend "▀", Mitte "█", unten Abschluss mit "▄"-Sequenzen.
        // Nur Spectre.Console Rendering APIs, keine direkten ANSI-Sequenzen.

        var version= "SEW                                v0.0.1";
        var top    = "█▀▀▀▀▀▀▀ ▄▀▀▀▀▀▄ █▀▀▀▀▀█ ▄▀▀▀▀▀▀ █▀▀▀▀▀▀▀";   // F     O     R     G     E
        var middle = "█▀▀▀▀▀▀  █     █ █▀▀▀▀█▀ █   ▀▀█ █▀▀▀▀▀▀ ";       // Vertikale Segmente / Öffnungen
        var bottom = "▀         ▀▀▀▀▀  ▀     ▀  ▀▀▀▀▀▀ ▀▀▀▀▀▀▀▀";     // Untere Abschlüsse / Innenräume

        var grid = new Grid();
        grid.AddColumn();
        grid.AddRow(new Markup("[bold red]" + version +"[/]"));
        grid.AddRow(new Markup(top));
        grid.AddRow(new Markup(middle));
        grid.AddRow(new Markup(bottom));
        grid.AddRow(new Markup("[green]Hallo Welt[/]"));
        return grid;
    }

    private static IRenderable BuildContent()
    {
        // Placeholder Content Panel – kann später durch dynamische Views ersetzt werden.
        var panel = new Panel(new Markup("[green]Inhalt[/] – Platz für zukünftige Module"))
        {
            Border = BoxBorder.Rounded,
            Padding = new Padding(1, 0, 1, 0)
        };
        panel.Header = new PanelHeader("Main");
        return panel;
    }

    private static IRenderable BuildFooter()
    {
        var grid = new Grid();
        grid.AddColumn();
        grid.AddRow(new Markup("[blue]Status:[/] Bereit"));
        grid.AddRow(new Markup("[dim]Forge © 2025[/]"));
        return grid;
    }
}
