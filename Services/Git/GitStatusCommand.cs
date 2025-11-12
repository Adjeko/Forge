using Spectre.Console;
using Spectre.Console.Rendering;

namespace Forge.Services.Git;

public sealed class GitStatusCommand : IGitCommand
{
    private readonly string _workingDirectory;
    public string Name => "status";

    public GitStatusCommand(string workingDirectory)
    {
        _workingDirectory = workingDirectory;
    }

    public IRenderable Execute()
    {
        try
        {
            var psi = new System.Diagnostics.ProcessStartInfo
            {
                FileName = "git",
                Arguments = "status --short --branch",
                WorkingDirectory = _workingDirectory,
                RedirectStandardOutput = true,
                RedirectStandardError = true,
                UseShellExecute = false,
                CreateNoWindow = true
            };
            using var proc = System.Diagnostics.Process.Start(psi)!;
            string stdout = proc.StandardOutput.ReadToEnd();
            string stderr = proc.StandardError.ReadToEnd();
            proc.WaitForExit();

            var content = string.IsNullOrWhiteSpace(stdout) ? "[grey]Keine Ausgabe[/]" : Escape(stdout);

            if (!string.IsNullOrWhiteSpace(stderr))
            {
                content += "\n[yellow]Warn/Error:[/] " + Escape(stderr);
            }

            return new Panel(new Markup(content))
            {
                Border = BoxBorder.Rounded,
                Header = new PanelHeader($"git status ({_workingDirectory})")
            };
        }
        catch (Exception ex)
        {
            return new Panel(new Markup($"[red]Fehler beim Ausführen von git status:[/] {Escape(ex.Message)}"))
            {
                Border = BoxBorder.Rounded,
                Header = new PanelHeader($"git status ({_workingDirectory})")
            };
        }
    }

    private static string Escape(string raw)
        => raw.Replace("[", "[[").Replace("]", "]]" );
}
