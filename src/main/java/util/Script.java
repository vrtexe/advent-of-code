package util;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.InputStreamReader;
import java.util.List;

import javax.script.Invocable;
import javax.script.ScriptEngine;
import javax.script.ScriptEngineFactory;
import javax.script.ScriptEngineManager;
import javax.script.ScriptException;

import com.google.gson.Gson;

public class Script {

  private final Invocable script;
  private final ScriptEngine engine;

  public Script(String filePath) {
    final var loadedScriptEngine = loadScript(filePath);
    this.engine = loadedScriptEngine;
    this.script = (Invocable) loadedScriptEngine;

  }

  public <T> T call(Class<T> clazz) {
    return this.script.getInterface(clazz);
  }

  public static void main(String[] args) {
    final var gson = new Gson();

    System.out.println(
        gson.toJson(new Hello("yoo1", "haha", 30)));

  }

  private record Hello(String yo, String lol, Integer yes) {
  }

  public static void listEngines() {
    ScriptEngineManager manager = new ScriptEngineManager(null);
    List<ScriptEngineFactory> engines = manager.getEngineFactories();

    for (ScriptEngineFactory engine : engines) {
      System.out.println("Engine name: {}".formatted(engine.getEngineName()));
      System.out.println("Version: {}".formatted(engine.getEngineVersion()));
      System.out.println("Language: {}".formatted(engine.getLanguageName()));

      System.out.println("Short Names:");
      for (String names : engine.getNames()) {
        System.out.println(names);
      }
    }
  }

  private static ScriptEngine loadScript(BufferedReader reader) {

    final var manager = new ScriptEngineManager(ClassLoader.getSystemClassLoader());
    final var engine = manager.getEngineByName("nashorn");

    try {
      engine.eval(reader);
    } catch (ScriptException e) {
      e.printStackTrace();
    }

    return engine;
  }

  private static ScriptEngine loadScript(String path) {
    final var file = new File(path);
    final var reader = createReader(file);

    return loadScript(reader);
  }

  private static BufferedReader createReader(File file) {
    try {
      return new BufferedReader(new InputStreamReader(new FileInputStream(file)));
    } catch (FileNotFoundException e) {
      e.printStackTrace();
    }

    return null;
  }

  public static double eval(final String str) {
    return new Object() {
      int pos = -1, ch;

      void nextChar() {
        ch = (++pos < str.length()) ? str.charAt(pos) : -1;
      }

      boolean eat(int charToEat) {
        while (ch == ' ')
          nextChar();
        if (ch == charToEat) {
          nextChar();
          return true;
        }
        return false;
      }

      double parse() {
        nextChar();
        double x = parseExpression();
        if (pos < str.length())
          throw new RuntimeException("Unexpected: " + (char) ch);
        return x;
      }

      // Grammar:
      // expression = term | expression `+` term | expression `-` term
      // term = factor | term `*` factor | term `/` factor
      // factor = `+` factor | `-` factor | `(` expression `)` | number
      // | functionName `(` expression `)` | functionName factor
      // | factor `^` factor

      double parseExpression() {
        double x = parseTerm();
        for (;;) {
          if (eat('+'))
            x += parseTerm(); // addition
          else if (eat('-'))
            x -= parseTerm(); // subtraction
          else
            return x;
        }
      }

      double parseTerm() {
        double x = parseFactor();
        for (;;) {
          if (eat('*'))
            x *= parseFactor(); // multiplication
          else if (eat('/'))
            x /= parseFactor(); // division
          else
            return x;
        }
      }

      double parseFactor() {
        if (eat('+'))
          return +parseFactor(); // unary plus
        if (eat('-'))
          return -parseFactor(); // unary minus

        double x;
        int startPos = this.pos;
        if (eat('(')) { // parentheses
          x = parseExpression();
          if (!eat(')'))
            throw new RuntimeException("Missing ')'");
        } else if ((ch >= '0' && ch <= '9') || ch == '.') { // numbers
          while ((ch >= '0' && ch <= '9') || ch == '.')
            nextChar();
          x = Double.parseDouble(str.substring(startPos, this.pos));
        } else if (ch >= 'a' && ch <= 'z') { // functions
          while (ch >= 'a' && ch <= 'z')
            nextChar();
          String func = str.substring(startPos, this.pos);
          if (eat('(')) {
            x = parseExpression();
            if (!eat(')'))
              throw new RuntimeException("Missing ')' after argument to " + func);
          } else {
            x = parseFactor();
          }
          if (func.equals("sqrt"))
            x = Math.sqrt(x);
          else if (func.equals("sin"))
            x = Math.sin(Math.toRadians(x));
          else if (func.equals("cos"))
            x = Math.cos(Math.toRadians(x));
          else if (func.equals("tan"))
            x = Math.tan(Math.toRadians(x));
          else
            throw new RuntimeException("Unknown function: " + func);
        } else {
          throw new RuntimeException("Unexpected: " + (char) ch);
        }

        if (eat('^'))
          x = Math.pow(x, parseFactor()); // exponentiation

        return x;
      }
    }.parse();
  }
}
