using Spectre.Console;
using Forge.Rendering;
using Spectre.Console.Rendering;

namespace Forge;

internal static class Program
{
	private static void Main()
	{
		AnsiConsole.Cursor.Hide();
		Render();

		int lastW = Console.WindowWidth;
		int lastH = Console.WindowHeight;
		while (true)
		{
			if (Console.KeyAvailable)
			{
				var key = Console.ReadKey(true);
				if (key.Key == ConsoleKey.Escape)
					break;

				// Ctrl+G -> git status im Zielverzeichnis ausführen
				if (key.Key == ConsoleKey.G && key.Modifiers.HasFlag(ConsoleModifiers.Control))
				{
					RunGitStatus(); // setzt _mainContent und rendert neu
					continue; // Nächste Iteration
				}
			}

			if (Console.WindowWidth != lastW || Console.WindowHeight != lastH)
			{
				lastW = Console.WindowWidth;
				lastH = Console.WindowHeight;
				Render();
			}

			Thread.Sleep(250);
		}
		AnsiConsole.Cursor.Show();
	}

	private static IRenderable? _mainContent;

	private static void Render()
	{
		AnsiConsole.Clear();
		// Beispiel: Farben können künftig konfigurierbar gemacht werden (CLI Argumente / Config).
		var layout = MainLayout.Build(startColor: "#FF0000", endColor: "#FFFFFF", mainContent: _mainContent);
		AnsiConsole.Write(layout);
	}

	private static void RunGitStatus()
	{
		try
		{
			// Keine Vorab-Statusmeldung: direkt Prozess starten
			var psi = new System.Diagnostics.ProcessStartInfo
			{
				FileName = "git",
				Arguments = "status --short --branch",
				WorkingDirectory = "C:/ADO/CS2",
				RedirectStandardOutput = true,
				RedirectStandardError = true,
				UseShellExecute = false,
				CreateNoWindow = true
			};
			using var proc = System.Diagnostics.Process.Start(psi)!;
			string stdout = proc.StandardOutput.ReadToEnd();
			string stderr = proc.StandardError.ReadToEnd();
			proc.WaitForExit();

			if (!string.IsNullOrWhiteSpace(stderr))
			{
				AnsiConsole.MarkupLine($"[yellow]Warn/Error:[/] {stderr.Replace("[", "[[").Replace("]", "]]")}");
			}

			if (string.IsNullOrWhiteSpace(stdout))
			{
				AnsiConsole.MarkupLine("[grey]Keine Ausgabe[/]");
				return;
			}

			// Einfache Rohtext-Ausgabe ohne Tabelle
			_mainContent = new Panel(new Markup(Escape(stdout)))
			{
				Border = BoxBorder.Rounded,
				Header = new PanelHeader("git status (C:/ADO/CS2)")
			};
			Render();
		}
		catch (Exception ex)
		{
			AnsiConsole.MarkupLine($"[red]Fehler beim Ausführen von git status:[/] {Escape(ex.Message)}");
		}
	}

	private static string Escape(string raw)
	{
		return raw.Replace("[", "[[").Replace("]", "]]");
	}
}
