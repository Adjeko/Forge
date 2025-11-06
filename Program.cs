using Spectre.Console;
using Forge.Rendering;

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

	private static void Render()
	{
		AnsiConsole.Clear();
		var layout = MainLayout.Build();
		AnsiConsole.Write(layout);
	}
}
