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
                new Layout("header").Size(5),
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
        var grid = new Grid();
        grid.AddColumn();
        grid.AddRow(new Markup("[bold yellow]Forge TUI[/]"));
        grid.AddRow(new Markup("[grey]Spectre.NET Layout Demo[/]"));
        grid.AddRow(new Markup($"Terminal: {Console.WindowWidth}x{Console.WindowHeight}"));
        grid.AddRow(new Markup($"Zeit: {DateTime.Now:HH:mm:ss}"));
        grid.AddRow(new Markup("[dim]ESC zum Beenden[/]"));
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
