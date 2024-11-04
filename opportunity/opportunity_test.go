package opportunity

import "testing"

func TestIsProfitableGain(t *testing.T) {
	op := &Opportunity{
		FirstCoinPrice:  30_000.0,
		SecondCoinPrice: 0.05,
		ThirdCoinPrice:  1_600.0,
		Fee:             0.001,
	}

	profitable, err := op.IsProfitable(1000)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	if !profitable {
		t.Errorf("This test should determine a gain, loss detected instead")
	}
}

func TestIsProfitableLoss(t *testing.T) {
	op := &Opportunity{
		FirstCoinPrice:  30_000.0,
		SecondCoinPrice: 0.05,
		ThirdCoinPrice:  1_500.0,
		Fee:             0.001,
	}

	profitable, err := op.IsProfitable(1000)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	if profitable {
		t.Errorf("This test should determine a loss, gain detected instead")
	}
}
