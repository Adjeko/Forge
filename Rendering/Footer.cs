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
        grid.AddRow(new Markup("[blue]Status:[/] Bereit"));
        grid.AddRow(new Markup("[dim]Forge © 2025[/]"));
        return grid;
    }
}