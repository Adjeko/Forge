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
					_mainContent = _gitStatusCommand.Execute();
					Render();
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
	private static Services.Git.IGitCommand _gitStatusCommand = new Services.Git.GitStatusCommand("C:/ADO/CS2");

	private static void Render()
	{
		AnsiConsole.Clear();
		// Beispiel: Farben können künftig konfigurierbar gemacht werden (CLI Argumente / Config).
		var layout = MainLayout.Build(startColor: "#FF0000", endColor: "#FFFFFF", mainContent: _mainContent);
		AnsiConsole.Write(layout);
	}




}
