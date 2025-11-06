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

        layout["header"].Update(new Header());
        layout["content"].Update(BuildContent());
    layout["footer"].Update(new Footer());

        return layout;
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

}
