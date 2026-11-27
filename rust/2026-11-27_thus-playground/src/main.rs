fn apply_discount(price: i32, discount: i32) -> i32 {
    price - price * discount / 100
}

#[test]
fn discount() {
    use thus::verify;

    let (price, discount, expected) = (1000, 10, 900);

    verify! {
        let actual = apply_discount(price, discount) => {
            actual == expected;
            actual >= 0 && actual <= price;
        },
        "price={price}, discount={discount}"
    }
}

fn main() {
    println!("10% OFF: {}", apply_discount(100, 10))
}
