using Spectre.Console;
using Spectre.Console.Rendering;

namespace Forge.Rendering;

/// <summary>
/// Separate Footer-Komponente als <see cref="IRenderable"/>. Zwei Zeilen: Status + Copyright.
/// Ausschließlich Spectre.Console APIs gemäß Richtlinien.
/// </summary>
internal sealed class Footer : IRenderable
{
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

    private static IRenderable BuildGrid()
    {
        var grid = new Grid();
        grid.AddColumn();

        // Hotkey-Legende (nur vorhandene, erweitert werden wenn neue Keys implementiert werden)
        var legendTable = new Table().NoBorder().AddColumn("Hotkeys");
        legendTable.HideHeaders();
		legendTable.AddRow(new Markup("[grey]Esc[/] : Beenden"));
		legendTable.AddRow(new Markup("[grey]Strg+G[/] : git status (C:/ADO/CS2)"));
        // Weitere Einträge können hier hinzugefügt werden, sobald zusätzliche Tasten unterstützt werden.
        grid.AddRow(legendTable);

        return grid;
    }
}