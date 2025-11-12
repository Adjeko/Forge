using Spectre.Console;
using Spectre.Console.Rendering;

namespace Forge.Rendering;

/// <summary>
/// Erzeugt das Grundlayout der Anwendung: 5 Zeilen Header, flexibler Inhalt, 2 Zeilen Footer.
/// Kapselt ausschließlich Spectre.Console Aufrufe gemäß Richtlinien.
/// </summary>
internal static class MainLayout
{
    /// <summary>
    /// Baut ein neues Layout für die aktuelle Render-Periode.
    /// </summary>
    public static Layout Build(string? startColor = null, string? endColor = null, IRenderable? mainContent = null)
    {
        var layout = new Layout("root")
            .SplitRows(
                new Layout("header").Size(5), // Header jetzt fünfzeilig für ASCII-FORGE
                new Layout("content"),
                new Layout("footer").Size(2)
            );

        // Parse Farbparameter falls vorhanden; erlaubt sowohl Hex (#RRGGBB) als auch Spectre Farbnamen.
        Color Parse(string? value, Color fallback)
        {
            if (string.IsNullOrWhiteSpace(value))
                return fallback;
            var v = value.Trim();
            if (v.StartsWith('#')) v = v[1..];
            if (v.Length == 6 && int.TryParse(v, System.Globalization.NumberStyles.HexNumber, null, out var hex))
            {
                var r = (hex >> 16) & 0xFF;
                var g = (hex >> 8) & 0xFF;
                var b = hex & 0xFF;
                return new Color((byte)r, (byte)g, (byte)b);
            }
            return v.ToLowerInvariant() switch
            {
                "red" => Color.Red,
                "white" => Color.White,
                "yellow" => Color.Yellow,
                "blue" => Color.Blue,
                "green" => Color.Green,
                "cyan" => new Color(0,255,255),
                "magenta" => new Color(255,0,255),
                "grey" => Color.Grey,
                _ => fallback
            };
        }

        var start = Parse(startColor, Color.Red);
        var end = Parse(endColor, Color.White);

        layout["header"].Update(new Header(start, end));
        layout["content"].Update(BuildContent(mainContent));
        layout["footer"].Update(new Footer());

        return layout;
    }

    private static IRenderable BuildContent(IRenderable? overrideContent)
    {
        var content = overrideContent ?? new Markup("[green]Inhalt[/] – Platz für zukünftige Module");
        // Wenn overrideContent bereits ein Panel ist, direkt zurückgeben um Doppelrahmen zu vermeiden
        if (overrideContent is Panel p)
            return p;

        var panel = new Panel(content)
        {
            Border = BoxBorder.Rounded,
            Padding = new Padding(0, 0, 0, 0)
        };
        return panel;
    }

}
