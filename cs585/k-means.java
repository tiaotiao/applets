
import java.util.*;

class Point {
    ArrayList<Double> w;
    private final int ATTRIBUTE_LEN = 2; 

    public Point() {
        this.w = new ArrayList<>();
        for (int i = 0; i < ATTRIBUTE_LEN; i++) {
            w.add(0.0);
        }
    }

    public Point(double w1, double w2) {
        this.w = new ArrayList<>();
        w.add(w1);
        w.add(w2);
    }

    public int size() {
        return w.size();
    }

    public double get(int index) {
        return w.get(index);
    }

    public void set(int index, double v) {
        w.set(index, v);
    }

    public double attrib(int index) {
        return w.get(index);
    }

    public boolean same(Point other) {
        for (int i = 0; i < w.size(); i++) {
            if (Math.abs(this.get(i) - other.get(i)) > 0.001) {
                return false;
            }
        }
        return true;
    }

    public double distance(Point other) {
        double d = 0;
        for (int i = 0; i < w.size(); i++) {
            double diff = (this.attrib(i) - other.attrib(i));
            d = d + diff * diff;
        }
        return Math.sqrt(d);
    }

    public String toString() {
        String s = "";
        int i = 0;
        for (double v: w) {
            if (i > 0) {
                s += ", ";
            }
            s += String.format("%.2f", v);
            i++;
        }
        return "(" + s + ")";
    }
}


class KMeans {

    private final int ATTRIBUTE_LEN = 2; 

    private ArrayList<List<Point>> grouping(ArrayList<Point> points, ArrayList<Point> centers) {
        int k = centers.size();

        ArrayList<List<Point>> groups = new ArrayList<>();
        for (int i = 0; i < k; i++) {
            groups.add(new ArrayList<Point>());
        }

        // group points
        for (Point p: points) {
            double minDist = -1;
            int centIndex = -1;

            // find the nearest center
            for (int i = 0; i < k; i++) {
                Point c = centers.get(i);
                double dist = p.distance(c);
                if (minDist < 0 || minDist > dist) {
                    minDist = dist;
                    centIndex = i;
                }
            }
            
            // put in group
            groups.get(centIndex).add(p);
        }

        return groups;
    }

    private ArrayList<Point> calcCenters(ArrayList<List<Point>> groups) {
        ArrayList<Point> centers = new ArrayList<>();
        for (List<Point> group: groups) {

            ArrayList<Double> sum = new ArrayList<>();
            for (int i = 0; i < ATTRIBUTE_LEN; i++) {
                sum.add(0.0);
            }

            for (Point p: group) {
                for (int i = 0; i < ATTRIBUTE_LEN; i++) {
                    sum.set(i, sum.get(i) + p.attrib(i));
                }
            }

            Point c = new Point();
            if (group.size() > 0) {
                for (int i = 0; i < ATTRIBUTE_LEN; i++) {
                    c.set(i, sum.get(i) / group.size());
                }
            }

            centers.add(c);
        }

        return centers;
    }

    private boolean sameCenters(ArrayList<Point> c1, ArrayList<Point> c2) {
        for (int i = 0; i < c1.size(); i++) {
            Point p1 = c1.get(i);
            Point p2 = c2.get(i);
            if (!p1.same(p2)) {
                return false;
            }
        }
        return true;
    }

    private void print(ArrayList<List<Point>> groups, ArrayList<Point> centers) {
        System.out.printf("Centers ---------------------------\n");
        for (int i = 0; i < centers.size(); i++) {
            Point c = centers.get(i);
            List<Point> group = groups.get(i);
            System.out.printf("Group %d:\t%s <<<\n", i, c.toString());
            for (Point p: group) {
                System.out.printf("\t\t%s\n", p);
            }
        }
    }

    public void k_means(ArrayList<Point> points, ArrayList<Point> centers) {
        
        ArrayList<List<Point>> groups;
        ArrayList<Point> newCenters;

        for (int i = 0; i < 3; i++) {
            groups = grouping(points, centers);
            newCenters = calcCenters(groups);


            if (sameCenters(centers, newCenters)) {
                break;
            }
            
            // output
            print(groups, newCenters);

            centers = newCenters;
        }

        System.out.println("Done.");
    }

}


class Main {
    
    public static void main(String[] args) {
        ArrayList<Point> points = new ArrayList<>();
        ArrayList<Point> centers = new ArrayList<>();
 
        /*
        points.add(new Point(2.0, 0.0));
        points.add(new Point(1.0, 3.0));
        points.add(new Point(3.0, 5.0));
        points.add(new Point(2.0, 2.0));
        points.add(new Point(4.0, 6.0));

        centers.add(new Point(1.0, 3.0));
        centers.add(new Point(2.0, 2.0));
        */

        points.add(new Point(1.0, 0.0));
        points.add(new Point(2.0, 4.0));
        points.add(new Point(0.0, 2.0));
        points.add(new Point(3.0, 5.0));
        points.add(new Point(1.0, 1.0));

        centers.add(new Point(0.0, 2.0));
        centers.add(new Point(1.0, 1.0));

        KMeans s = new KMeans();
        s.k_means(points, centers);        
    }
}