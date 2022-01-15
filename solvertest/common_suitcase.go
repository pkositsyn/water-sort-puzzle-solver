package solvertest

import (
	"fmt"
	"log"
	"strings"
	"testing"

	watersortpuzzle "github.com/pkositsyn/water-sort-puzzle-solver"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func (s *SolverSuite) TestSolver() {
	testCases := []struct {
		state         string
		expectedSteps int
	}{
		{
			state:         "O;OOO",
			expectedSteps: 1,
		},
		{
			state:         "FOFO;OFOF;",
			expectedSteps: 7,
		},
		{
			state:         "FORF;OORF;RFOR;;",
			expectedSteps: 10,
		},
		{
			state:         "FROO;FRFR;OFRO;;",
			expectedSteps: 10,
		},
		{
			state:         "RGGG;ORPG;PORO;FPOP;FFFR;;",
			expectedSteps: 12,
		},
		{
			state:         "GORO;FFRO;PPFO;GPRF;GRGP;;",
			expectedSteps: 15,
		},
		{
			state:         "GOGF;OPPO;PRFR;FRGP;FGRO;;",
			expectedSteps: 16,
		},
		{
			state:         "PRFP;RGGO;ROOP;PRGF;GOFF;;",
			expectedSteps: 14,
		},
		{
			state:         "FPFB;PPGB;OOQO;BRPO;FGRQ;QFRR;QBGG;;",
			expectedSteps: 20,
		},
		{
			state:         "FPGG;OFPF;FORG;OGRP;RORP;;",
			expectedSteps: 16,
		},
		{
			state:         "FPGR;OGGB;PQOR;GRFB;BPQB;POFQ;QRFO;;",
			expectedSteps: 22,
		},
		{
			state:         "QFFF;FQPO;QOQG;ROGP;RBPR;OBGB;PGBR;;",
			expectedSteps: 21,
		},
		{
			state:         "RPPR;BRFR;OGOO;QFGQ;GBQO;BQPB;PFFG;;",
			expectedSteps: 21,
		},
		{
			state:         "BQFB;PFRG;FPGF;BRQO;GOBG;RPOR;OPQQ;;",
			expectedSteps: 22,
		},
		{
			state:         "GOFR;OPRG;OFRG;PFFG;PPRO;;",
			expectedSteps: 16,
		},
		{
			state:         "GRPP;GBPB;FOQQ;OPGQ;FGBR;FFBQ;OORR;;",
			expectedSteps: 20,
		},
		{
			state:         "ORRF;PGRO;FFGR;GOPF;OPGP;;",
			expectedSteps: 15,
		},
		{
			state:         "BBFG;QROP;RGOF;QFRP;QOPP;GBFB;GQRO;;",
			expectedSteps: 22,
		},
		{
			state:         "GRPO;PRFB;OQOB;PGFB;PQRB;QGGR;FFOQ;;",
			expectedSteps: 21,
		},
		{
			state:         "RFFF;GGOO;GRPO;RGOP;PRPF;;",
			expectedSteps: 13,
		},
		{
			state:         "ORBB;GPPG;QFOG;PFQR;OQPG;RROB;BFFQ;;",
			expectedSteps: 19,
		},
		{
			state:         "OORP;RGGF;ORPP;PFFG;FRGO;;",
			expectedSteps: 13,
		},
		{
			state:         "OQBF;PPRP;OQGQ;GFPR;FFBQ;ROOB;BGGR;;",
			expectedSteps: 19,
		},
		{
			state:         "QBPO;BGGP;FOFO;PBGF;QRGF;BQQR;RORP;;",
			expectedSteps: 22,
		},
		{
			state:         "FGPT;BTHF;FQGO;POOB;QRRP;FOHG;GRTB;QHRH;PBQT;;",
			expectedSteps: 29,
		},
		{
			state:         "GOPO;OFTQ;TQRP;BHQR;GFRH;QPHR;BGOG;FBBT;HTPF;;",
			expectedSteps: 28,
		},
		{
			state:         "RGGR;BFOP;QQPF;BGBO;GOBF;PQQR;PFOR;;",
			expectedSteps: 21,
		},
		{
			state:         "TRFH;QFOO;QGQG;THBT;BRRB;FPQP;ORPF;OPBH;HGTG;;",
			expectedSteps: 28,
		},
		{
			state:         "BRQF;GRFG;GFBP;RRGP;QBOP;QOPB;OOFQ;;",
			expectedSteps: 21,
		},
		{
			state:         "QGFH;QTGG;OQRP;BBTH;HFRB;RFOR;PTBQ;POOH;GPTF;;",
			expectedSteps: 27,
		},
		{
			state:         "FBHB;FRHT;QTFF;RPOG;QGPR;OGGH;HQTR;TQPO;OBBP;;",
			expectedSteps: 27,
		},
		{
			state:         "FBPG;ROQP;BFFO;POBR;PFGO;QGBR;GQRQ;;",
			expectedSteps: 22,
		},
		{
			state:         "OHTP;TGFR;FGBF;ORRB;PTQT;HFBQ;QOHG;RPHP;BGQO;;",
			expectedSteps: 28,
		},
		{
			state:         "POGR;OBFP;OQGP;PQRG;BQQB;GRFO;FBRF;;",
			expectedSteps: 22,
		},
		{
			state:         "ROGB;PTQH;BQGP;HOOG;ROTR;PFTT;HBFQ;FRBP;QFHG;;",
			expectedSteps: 28,
		},
		{
			state:         "GPBH;PBBF;RORR;QTQH;BPHG;TTOG;ROQH;GFFO;FPTQ;;",
			expectedSteps: 26,
		},
		{
			state:         "QTGO;HBFQ;OHFB;GQHR;TGTP;PBRT;PBRP;GQOH;ROFF;;",
			expectedSteps: 28,
		},
		{
			state:         "GTFH;HORF;BHQP;PGRO;BBGO;QQOT;GFFR;HBPT;QRTP;;",
			expectedSteps: 28,
		},
		{
			state:         "BRGO;BGFF;QORB;GPPF;PRQO;RQBO;FGQP;;",
			expectedSteps: 21,
		},
		{
			state:         "GBRP;FFFG;FROG;OROB;BPQQ;HHTO;QHGB;HPPT;TTRQ;;",
			expectedSteps: 24,
		},
	}

	for i, testCase := range testCases {
		tt := testCase

		s.Run(fmt.Sprintf("Test %d", i), func() {
			testFlasks := len(strings.Split(tt.state, ";"))

			if s.MaxFlasks != 0 && s.MaxFlasks < testFlasks {
				s.T().Skipf("Test with %d flasks skipped, MaxFlasks is %d", testFlasks, s.MaxFlasks)
			}

			solver := s.NewSolverFunc()

			var initialState watersortpuzzle.State
			s.Require().NoError(initialState.FromString(tt.state))

			steps, err := solver.Solve(initialState)
			s.Require().NoError(err)

			state := initialState
			for _, step := range steps {
				state, err = state.Step(step)
				s.Require().NoError(err)
			}

			if aStarSolver, ok := solver.(*watersortpuzzle.AStarSolver); ok {
				log.Printf("%+v, Path length: %d\n", aStarSolver.Stats(), len(steps))
			}
			s.Assert().Equal(tt.expectedSteps, len(steps))
			s.Assert().True(state.IsTerminal())
		})
	}
}

func TemplateBenchmarkSolve(b *testing.B, newSolverFunc func() watersortpuzzle.Solver) {
	const state = "ORRF;PGRO;FFGR;GOPF;OPGP;;"

	var initialState watersortpuzzle.State
	require.NoError(b, initialState.FromString(state))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		solver := newSolverFunc()
		_, err := solver.Solve(initialState)
		require.NoError(b, err)
	}
}

type SolverSuite struct {
	suite.Suite
	NewSolverFunc func() watersortpuzzle.Solver
	MaxFlasks     int
}
