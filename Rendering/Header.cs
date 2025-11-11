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

    private readonly Color _startColor;
    private readonly Color _endColor;

    public Header(Color? start = null, Color? end = null)
    {
        _startColor = start ?? Color.Red;
        _endColor = end ?? Color.White;
    }

    private IRenderable BuildGrid()
    {
        int width = AnsiConsole.Profile.Width;
        const string LeftPattern = "////";
        const string FillerPattern = "////";

        // Sichtbare Länge eines Markup-Strings (ignoriert Spectre-Markup Tags)
        static int GetVisibleLength(string content)
        {
            int len = 0;
            bool inTag = false;
            for (int i = 0; i < content.Length; i++)
            {
                var c = content[i];
                if (c == '[')
                {
                    // Start eines Tags
                    inTag = true;
                    continue;
                }
                if (inTag && c == ']')
                {
                    inTag = false;
                    continue;
                }
                if (!inTag)
                    len++;
            }
            return len;
        }

        // Erzeugt eine Zeile mit optionaler einfacher Farb-Markup Kennzeichnung
        string BuildLine(string content, string? color = null)
        {
            var main = color is null ? content : $"[{color}]{content}[/]";
            int visibleLen = LeftPattern.Length + 1 + GetVisibleLength(content) + 1;
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

        // Gradient-Funktion: Interpoliert von Rot (255,0,0) nach Weiß (255,255,255)
        string ApplyGradient(string text, Color startColor, Color endColor)
        {
            // Farbstufen nur über nicht-leere Zeichen anwenden; Spaces bleiben unverändert
            var chars = text.ToCharArray();
            int gradientChars = chars.Count(c => c != ' ');
            if (gradientChars <= 1)
                return text; // Nichts zu graduieren

            int processed = 0;
            var sb = new System.Text.StringBuilder(text.Length * 20);
            foreach (var ch in chars)
            {
                if (ch == ' ')
                {
                    sb.Append(ch);
                    continue;
                }
                double t = processed / (double)(gradientChars - 1); // 0..1
                // Lineare Interpolation zwischen Start- und Endfarbe
                int r = (int)(startColor.R + ((endColor.R - startColor.R) * t));
                int g = (int)(startColor.G + ((endColor.G - startColor.G) * t));
                int b = (int)(startColor.B + ((endColor.B - startColor.B) * t));
                sb.Append("[#");
                sb.Append(r.ToString("X2"));
                sb.Append(g.ToString("X2"));
                sb.Append(b.ToString("X2"));
                sb.Append(']');
                sb.Append(ch);
                sb.Append("[/]");
                processed++;
            }
            return sb.ToString();
        }

        var lines = new[]
        {
            BuildLine(GetVersionText(), "bold red"),
            BuildLine(ApplyGradient(Top, _startColor, _endColor)),
            BuildLine(ApplyGradient(Middle, _startColor, _endColor)),
            BuildLine(ApplyGradient(Bottom, _startColor, _endColor)),
            BuildLine("/////////////////////////////////////////", "red")
        };

        var grid = new Grid();
        grid.AddColumn(new GridColumn().Padding(0,0,0,0));
        foreach (var line in lines)
            grid.AddRow(new Markup(line));
        return grid;
    }
}
