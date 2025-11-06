using Spectre.Console;
using System.Text;

namespace Forge;

internal static class Program
{
	private static void Main()
	{
		Console.CursorVisible = false;
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

			Thread.Sleep(150);
		}
		Console.CursorVisible = true;
	}

	private static void Render()
	{
	}
}
