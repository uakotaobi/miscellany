// -*- compile-command: "javac -d bin -sourcepath . -classpath '.' Application.java && java -classpath bin org.team1759.Application" -*-

package org.team1759;
import java.util.List;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.Random;
import java.io.FileWriter;
import java.io.IOException;

public class Application {

    private static final int width = 500;
    private static final int height = 500;

    private static class RGB {
        private double r_ = 0, g_ = 0, b_ = 0; // 0 <= r, g, b <= 255

        public RGB() { /* All channels already 0.  Nothing to do. */ }
        public RGB(double red, double green, double blue) {
            r_ = red;
            g_ = green;
            b_ = blue;
        }
        public double r() { return r_; }
        public double g() { return g_; }
        public double b() { return b_; }
        public String toString() {
            return String.format("#%02x%02x%02x", (int)r_, (int)g_, (int)b_);
        }
    }

    private static class Point {
        public double x = 0, y = 0;
        public Point() { }
        public Point(double x_, double y_) { x = x_; y = y_; }
    }

    private static List<RGB> pixels = new ArrayList<RGB>(Collections.nCopies(width * height, new RGB(0, 0, 0)));

    public static void main(String[] args) {

        Random rand = new Random();
        // List<Point> points = Arrays.asList(new Point(rand.nextInt(width), rand.nextInt(height)),
        //                                    new Point(rand.nextInt(width), rand.nextInt(height)),
        //                                    new Point(rand.nextInt(width), rand.nextInt(height))
        // );
        List<Point> points = Arrays.asList(new Point(width/2, 0),
                                           new Point(0, height - 1),
                                           new Point(width - 1, height - 1));
        chaosGame(points, 100000);
        createImage("foo.ppm");
    }

    private static void setPixel(int x, int y, RGB color) {
        int offset = width * y + x;
        pixels.set(offset, color);
    }

    /// Writes an ASCII Portable Pixmap (PPM) file containing the image
    /// data in the pixels array.
    private static void createImage(String filename) {

        try (FileWriter writer = new FileWriter(filename)) {
            writer.write("P3\n# " + filename + "\n");
            writer.write(String.valueOf(width) + " ");
            writer.write(String.valueOf(height) + "\n");
            writer.write("255\n");

            for (RGB color : pixels) {
                // To make sure "no line is longer than 70 characters,"
                // put a newline at the end of every pixel.
                writer.write(String.valueOf((int)color.r()) + " ");
                writer.write(String.valueOf((int)color.g()) + " ");
                writer.write(String.valueOf((int)color.b()) + "\n");
            }

        } catch (IOException e) {
            System.out.println("The program blew up: " + e.getMessage());
        } finally {
            System.out.println("Finally block executed.");
        }
    }

    private static void chaosGame(List<Point> seeds, int iterations) {
        Random rand = new Random();
        List<RGB> colors = Arrays.asList(new RGB(255, 0, 0), new RGB(0, 255, 0), new RGB(0, 0, 255));

        // Choose a random initial point.
        Point initial = seeds.get(rand.nextInt(seeds.size()));
        Point current = new Point(initial.x, initial.y);

        for (int i = 0; i < iterations; ++i) {
            // Determine the current color based on weight averages of the
            // vertex colors.
            double[] weights = computeBarycentricCoordinates(seeds.get(0),
                                                             seeds.get(1),
                                                             seeds.get(2),
                                                             current);
            RGB currentColor = new RGB(weights[0] * colors.get(0).r() + weights[1] * colors.get(1).r() + weights[2] * colors.get(2).r(),
                                       weights[0] * colors.get(0).g() + weights[1] * colors.get(1).g() + weights[2] * colors.get(2).g(),
                                       weights[0] * colors.get(0).b() + weights[1] * colors.get(1).b() + weights[2] * colors.get(2).b());

            // Draw the current point.
            if (current.x >= 0 && current.x < width && current.y >= 0 && current.y < height) {
                setPixel((int)current.x, (int)current.y, currentColor);
            }

            // Choose random target point.
            Point next = seeds.get(rand.nextInt(seeds.size()));

            // Move to the point halfway between here and there.
            current.x = (current.x + next.x) / 2;
            current.y = (current.y + next.y) / 2;
        }
    }

    // What are barycentric coordinates?  They're a way of providing a
    // normalized coordinate system over an N-simplex polytope, in the same
    // way that people use (u, v) pairs to form a coordinate system over
    // parametric surfaces.
    //
    // - For a given N-simplex (lines, triangles, tetrahedra, 5-cells, and so
    //   on), there are N + 1 weights in the barycentric coordinate.
    // - These weights always add to 1.0.
    // - By definition, if only one of the weights is nonzero, then the
    //   barycentric coordinate lies on a vertex of the N-simplex.
    // - By definition, if only two of the weights are nonzero, then the
    //   barycentric coordinate lies along an edge of the N-simplex.
    //
    // By using these coordinates as weighted averages over a set of colors,
    // we ensure an even interpolation of those colors across the simplex.
    //
    // This function only deals with the 2-simplex (a triangle), so it returns
    // an array of three weights (one per vertex.)  P is the point you're
    // trying to get the barycentric coordinates for, while a, b, and c are
    // the vertices of the triangle.
    //
    // Let's be frank: I stole this formula whole-hog from
    // https://en.wikipedia.org/wiki/Barycentric_coordinate_system#Conversion_between_barycentric_and_Cartesian_coordinates.
    static double[] computeBarycentricCoordinates(Point a, Point b, Point c, Point p) {
        double determinant = (b.y - c.y) * (a.x - c.x) + (c.x - b.x) * (a.y - c.y);
        double l1 = ((b.y - c.y) * (p.x - c.x) + (c.x - b.x) * (p.y - c.y)) / determinant;
        double l2 = ((c.y - a.y) * (p.x - c.x) + (a.x - c.x) * (p.y - c.y)) / determinant;
        double l3 = 1 - l1 - l2;
        return new double[] { l1, l2, l3 };
    }
}
