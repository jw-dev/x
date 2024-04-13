package pgn

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractMoves(t *testing.T) {
	tests := []struct {
		m   string
		exp int
	}{
		{
			m:   "1. e4 { [%eval 0.17] } 1... c5 { [%eval 0.28] } 2. Bc4 { [%eval 0.11] } 2... Nc6 { [%eval 0.0] } 3. Qh5? { [%eval -1.14] } 3... g6? { [%eval 0.22] } 4. Qxc5 { [%eval 0.17] } 4... d6 { [%eval 0.58] } 5. Qe3 { [%eval 0.53] } 5... e5 { [%eval 0.75] } 6. Bd5?! { [%eval 0.14] } 6... Nd4 { [%eval -0.03] } 7. Qc3?! { [%eval -0.98] } 7... Nf6 { [%eval -0.63] } 8. Bc4?? { [%eval -5.0] } 8... Nxe4 { [%eval -5.03] } 9. Qd3 { [%eval -4.85] } 9... Nc5 { [%eval -4.96] } 10. Qc3 { [%eval -5.12] } 10... Be7? { [%eval -2.73] } 11. b4? { [%eval -4.92] } 11... Ne4 { [%eval -4.91] } 12. Qe3?? { [%eval -10.55] } 12... Nxc2+ { [%eval -10.42] } 13. Ke2 { [%eval -11.0] } 13... Nxe3 { [%eval -11.12] } 14. Kxe3 { [%eval -18.54] } 14... Ng5 { [%eval -10.74] } 15. Nf3 { [%eval -14.2] } 15... Nxf3 { [%eval -12.66] } 16. Kxf3 { [%eval -14.48] } 16... O-O { [%eval -10.81] } 17. Re1 { [%eval -13.67] } 17... h5 { [%eval -9.77] } 18. Ke2 { [%eval -10.91] } 18... Bg4+?! { [%eval -9.47] } 19. Kf1 { [%eval -9.38] } 19... Bg5?! { [%eval -8.71] } 20. h3 { [%eval -9.06] } 20... Be6 { [%eval -9.43] } 21. Ba3?! { [%eval -13.42] } 21... Bxc4+ { [%eval -13.63] } 22. Kg1 { [%eval -13.57] } 22... Bxd2? { [%eval -8.57] } 23. Nxd2 { [%eval -8.44] } 23... Bb5 { [%eval -8.13] } 24. Ne4?! { [%eval -8.63] } 24... d5 { [%eval -8.54] } 25. Nc5?! { [%eval -9.13] } 25... b6 { [%eval -8.97] } 26. Nb3 { [%eval -8.9] } 26... d4?! { [%eval -8.15] } 27. Nd2?! { [%eval -9.07] } 27... Re8 { [%eval -9.24] } 28. Ne4?! { [%eval -10.56] } 28... d3? { [%eval -8.62] } 29. Bc1 { [%eval -8.64] } 29... f5 { [%eval -8.15] } 30. Nd2? { [%eval -10.1] } 30... e4?! { [%eval -9.16] } 31. a4?! { [%eval -9.84] } 31... Ba6?! { [%eval -9.24] } 32. b5?! { [%eval -11.08] } 32... Bb7 { [%eval -11.05] } 33. Bb2 { [%eval -12.79] } 33... Rc8 { [%eval -11.79] } 34. Rac1 { [%eval -17.82] } 34... Rxc1 { [%eval -13.43] } 35. Rxc1 { [%eval -16.56] } 35... Qd7 { [%eval -11.55] } 36. f3 { [%eval -22.54] } 36... exf3 { [%eval -20.36] } 37. gxf3 { [%eval -30.72] } 37... Rc8 { [%eval -16.03] } 38. Re1 { [%eval -21.03] } 38... Rc2 { [%eval -19.66] } 39. Bc1?! { [%eval #-11] } 39... Rxd2?! { [%eval -21.07] } 40. Bxd2 { [%eval -18.86] } 40... Qd4+ { [%eval -13.44] } 41. Kh2?! { [%eval #-2] } 41... Bxf3?! { [%eval #-10] } 42. Re8+ { [%eval #-8] } 42... Kf7 { [%eval #-7] } 43. Rc8 { [%eval #-1] } 43... Qf2# 0-1",
			exp: 86,
		},
		{
			m:   "1. e4 e6 2. d4 b6 3. a3 Bb7 4. Nc3 Nh6 5. Bxh6 gxh6 6. Be2 Qg5 7. Bg4 h5 8. Nf3 Qg6 9. Nh4 Qg5 10. Bxh5 Qxh4 11. Qf3 Kd8 12. Qxf7 Nc6 13. Qe8# 1-0",
			exp: 25,
		},
		{
			m:   "1/2-1/2",
			exp: 0,
		},
		{
			m:   "1. e4 1-0",
			exp: 1,
		},
	}
	r := Result{}
	for _, test := range tests {
		r.extractMoves(test.m)
		require.Equal(t, test.exp, len(r.Moves))
	}
}
