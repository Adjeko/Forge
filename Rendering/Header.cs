using Spectre.Console;
using Spectre.Console.Rendering;
using System.Reflection;

namespace Forge.Rendering;

/// <summary>
/// Separate Header-Komponente als <see cref="IRenderable"/>.
/// Erzeugt fünf Zeilen: Version, ASCII-Schriftzug (3 Zeilen), zusätzliche Status-/Grußzeile.
/// Ausschließlich Spectre.Console APIs gemäß Richtlinien.
/// </summary>
internal sealed class Header : IRenderable
{
    // Einmal berechnete formatierte Versionszeile
    private static readonly string VersionLine = BuildVersionLine();
    private static string GetVersionText() => VersionLine;

    // ASCII Darstellung – bewusst in Konstanten gelassen für einfache Anpassung
    private const string Top    = "█▀▀▀▀▀▀▀ ▄▀▀▀▀▀▄ █▀▀▀▀▀█ ▄▀▀▀▀▀▀ █▀▀▀▀▀▀▀";   // F O R G E
    private const string Middle = "█▀▀▀▀▀▀  █     █ █▀▀▀▀█▀ █   ▀▀█ █▀▀▀▀▀▀ ";
    private const string Bottom = "▀         ▀▀▀▀▀  ▀     ▀  ▀▀▀▀▀▀ ▀▀▀▀▀▀▀▀";

    public IEnumerable<Segment> Render(RenderOptions options, int maxWidth)
    {
        var grid = BuildGrid();
        foreach (var segment in grid.Render(options, maxWidth))
            yield return segment;
    }

    public Measurement Measure(RenderOptions options, int maxWidth)
    {
        var grid = BuildGrid();
        return grid.Measure(options, maxWidth);
    }

    private static string BuildVersionLine()
    {
        var asm = typeof(Program).Assembly;
        var informational = asm.GetCustomAttribute<AssemblyInformationalVersionAttribute>()?.InformationalVersion;
        var file = asm.GetCustomAttribute<AssemblyFileVersionAttribute>()?.Version;
        var nameVersion = asm.GetName().Version?.ToString();
        var raw = informational ?? file ?? nameVersion ?? "0.0.0";

        var cut = raw.Split('-', '+')[0];
        var parts = cut.Split('.');
        if (parts.Length >= 3)
            cut = string.Join('.', parts[0], parts[1], parts[2]);
        else if (parts.Length == 2)
            cut = string.Join('.', parts[0], parts[1], "0");
        else if (parts.Length == 1)
            cut = parts[0] + ".0.0";

        var version = cut.StartsWith("v") ? cut : "v" + cut;

        const string prefix = "SEW";
        int totalWidth = Top.Length;
        int availableForVersion = totalWidth - prefix.Length;
        if (availableForVersion <= 0)
            return prefix;

        var versionDisplay = version.Length > availableForVersion
            ? version.Substring(0, availableForVersion)
            : version;

        int spacingLen = totalWidth - (prefix.Length + versionDisplay.Length);
        if (spacingLen < 0) spacingLen = 0;
        return prefix + new string(' ', spacingLen) + versionDisplay;
    }

    private IRenderable BuildGrid()
    {
        int width = AnsiConsole.Profile.Width;
        const string LeftPattern = "////";
        const string FillerPattern = "////";

        string BuildLine(string content, string? color = null)
        {
            var main = color is null ? content : $"[{color}]{content}[/]";
            int visibleLen = LeftPattern.Length + 1 + content.Length + 1;
            int remaining = width - visibleLen;
            if (remaining <= 0)
                return $"[red]{LeftPattern}[/] {main}";

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
            BuildLine(GetVersionText(), "bold red"),
            BuildLine(Top),
            BuildLine(Middle),
            BuildLine(Bottom),
            BuildLine("/////////////////////////////////////////", "red")
        };

        var grid = new Grid();
        grid.AddColumn(new GridColumn().Padding(0,0,0,0));
        foreach (var line in lines)
            grid.AddRow(new Markup(line));
        return grid;
    }
}
