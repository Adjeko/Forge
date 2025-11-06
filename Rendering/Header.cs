using Spectre.Console;
using Spectre.Console.Rendering;

namespace Forge.Rendering;

/// <summary>
/// Separate Header-Komponente als <see cref="IRenderable"/>.
/// Erzeugt fünf Zeilen: Version, ASCII-Schriftzug (3 Zeilen), zusätzliche Status-/Grußzeile.
/// Ausschließlich Spectre.Console APIs gemäß Richtlinien.
/// </summary>
internal sealed class Header : IRenderable
{
    private const string VersionText = "SEW                                v0.0.1"; // TODO: Dynamisch versionieren

    // ASCII Darstellung – bewusst in Konstanten gelassen für einfache Anpassung
    private const string Top    = "█▀▀▀▀▀▀▀ ▄▀▀▀▀▀▄ █▀▀▀▀▀█ ▄▀▀▀▀▀▀ █▀▀▀▀▀▀▀";   // F O R G E
    private const string Middle = "█▀▀▀▀▀▀  █     █ █▀▀▀▀█▀ █   ▀▀█ █▀▀▀▀▀▀ ";
    private const string Bottom = "▀         ▀▀▀▀▀  ▀     ▀  ▀▀▀▀▀▀ ▀▀▀▀▀▀▀▀";

    /// <summary>
    /// Erzeugt alle Segmente für den Header.
    /// </summary>
    public IEnumerable<Segment> Render(RenderOptions options, int maxWidth)
    {
        var grid = BuildGrid();
        foreach (var segment in grid.Render(options, maxWidth))
            yield return segment;
    }

    /// <summary>
    /// Misst die Breite/Höhe des Grids (delegiert an internes Grid).
    /// </summary>
    public Measurement Measure(RenderOptions options, int maxWidth)
    {
        var grid = BuildGrid();
        return grid.Measure(options, maxWidth);
    }

    private IRenderable BuildGrid()
    {
        int width = AnsiConsole.Profile.Width; // Aktuelle Konsolenbreite
        const string LeftPattern = "////";
        const string FillerPattern = "////";

        string BuildLine(string content, string? color = null)
        {
            var main = color is null ? content : $"[{color}]{content}[/]";
            int visibleLen = LeftPattern.Length + 1 + content.Length + 1; // Pattern + Space + content + Space
            int remaining = width - visibleLen;
            if (remaining <= 0)
                return $"[red]{LeftPattern}[/] {main}"; // Kein Platz für rechten Füller

            var fillerBuilder = new System.Text.StringBuilder(remaining);
            while (fillerBuilder.Length < remaining)
            {
                var next = FillerPattern;
                if (fillerBuilder.Length + next.Length > remaining)
                    next = next.Substring(0, remaining - fillerBuilder.Length);
                fillerBuilder.Append(next);
            }
            var filler = fillerBuilder.ToString();
            return $"[red]{LeftPattern}[/] {main} [red]{filler}[/]";
        }

        var lines = new[]
        {
            BuildLine(VersionText, "bold red"),
            BuildLine(Top),
            BuildLine(Middle),
            BuildLine(Bottom),
            BuildLine("Hallo Welt", "green")
        };

        var grid = new Grid();
        grid.AddColumn(new GridColumn().Padding(0,0,0,0));
        foreach (var line in lines)
            grid.AddRow(new Markup(line));
        return grid;
    }
}
