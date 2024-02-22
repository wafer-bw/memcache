package expire_test

//go:generate go run go.uber.org/mock/mockgen@latest -source=expire.go -destination=../mocks/mockexpire/mockexpire.go -package=mockexpire

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache/internal/expire"
	"github.com/wafer-bw/memcache/internal/mocks/mockexpire"
	"github.com/wafer-bw/memcache/internal/ports"
	"go.uber.org/mock/gomock"
)

var _ ports.Expirer[int, int] = (*expire.AllKeys[int, int])(nil)

func TestAllKeys_Expire(t *testing.T) {
	t.Parallel()

	t.Run("deletes all expired items", func(t *testing.T) {
		t.Parallel()

		expired := time.Until(time.Now().Add(-1 * time.Minute))
		unexpired := time.Until(time.Now().Add(1 * time.Minute))
		ctrl := gomock.NewController(t)
		m := mockexpire.NewMockCacher[int, int](ctrl)
		sut := expire.AllKeys[int, int]{}

		gomock.InOrder(
			m.EXPECT().Keys().Return([]int{1, 2, 3, 4, 5}),
			m.EXPECT().TTL(1).Return(&expired, true),
			m.EXPECT().Delete(1),
			m.EXPECT().TTL(2).Return(&unexpired, true),
			m.EXPECT().TTL(3).Return(nil, false),
			m.EXPECT().TTL(4).Return(&expired, true),
			m.EXPECT().Delete(4),
			m.EXPECT().TTL(5).Return(&expired, true),
			m.EXPECT().Delete(5),
		)

		sut.Expire(m)
	})
}

func TestRandomSample_Expire(t *testing.T) {
	t.Parallel()

	t.Run("expires items under percentage and stops", func(t *testing.T) {
		t.Parallel()

		expired := time.Until(time.Now().Add(-1 * time.Minute))
		unexpired := time.Until(time.Now().Add(1 * time.Minute))
		ctrl := gomock.NewController(t)
		m := mockexpire.NewMockCacher[int, int](ctrl)
		sut := expire.RandomSample[int, int]{
			SampleSize:    20,
			ExpirePercent: 0.25,
		}

		keys := make([]int, 100)
		for i := 0; i < 100; i++ {
			keys[i] = i
		}

		expects := make([]any, 0, 25)
		expects = append(expects, m.EXPECT().Size().Return(100))
		for i := 10; i < 15; i++ {
			expects = append(expects,
				m.EXPECT().RandomKey().Return(i, true),
				m.EXPECT().TTL(i).Return(&expired, true),
				m.EXPECT().Delete(i),
			)
		}
		for i := 15; i < 30; i++ {
			expects = append(expects,
				m.EXPECT().RandomKey().Return(i, true),
				m.EXPECT().TTL(i).Return(&unexpired, true),
			)
		}

		gomock.InOrder(expects...)

		sut.Expire(m)
	})

	t.Run("expires items until under percentage then stops", func(t *testing.T) {
		t.Parallel()

		expired := time.Until(time.Now().Add(-1 * time.Minute))
		unexpired := time.Until(time.Now().Add(1 * time.Minute))
		ctrl := gomock.NewController(t)
		m := mockexpire.NewMockCacher[int, int](ctrl)
		sut := expire.RandomSample[int, int]{
			SampleSize:    20,
			ExpirePercent: 0.25,
		}

		keys := make([]int, 100)
		for i := 0; i < 100; i++ {
			keys[i] = i
		}

		expects := make([]any, 0, 35)

		// Round 1
		expects = append(expects, m.EXPECT().Size().Return(100))
		for i := 10; i < 20; i++ {
			expects = append(expects,
				m.EXPECT().RandomKey().Return(i, true),
				m.EXPECT().TTL(i).Return(&expired, true),
				m.EXPECT().Delete(i),
			)
		}
		for i := 20; i < 30; i++ {
			expects = append(expects,
				m.EXPECT().RandomKey().Return(i, true),
				m.EXPECT().TTL(i).Return(&unexpired, true),
			)
		}

		// Round 2
		expects = append(expects, m.EXPECT().Size().Return(90))
		for i := 30; i < 35; i++ {
			expects = append(expects,
				m.EXPECT().RandomKey().Return(i, true),
				m.EXPECT().TTL(i).Return(&expired, true),
				m.EXPECT().Delete(i),
			)
		}
		for i := 35; i < 50; i++ {
			expects = append(expects,
				m.EXPECT().RandomKey().Return(i, true),
				m.EXPECT().TTL(i).Return(&unexpired, true),
			)
		}

		gomock.InOrder(expects...)

		sut.Expire(m)
	})

	t.Run("returns early if size is found to be 0", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		m := mockexpire.NewMockCacher[int, int](ctrl)
		sut := expire.RandomSample[int, int]{SampleSize: 20, ExpirePercent: 0.25}

		m.EXPECT().Size().Return(0)

		sut.Expire(m)
	})

	t.Run("returns early if sampling a random key finds no key", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		m := mockexpire.NewMockCacher[int, int](ctrl)
		sut := expire.RandomSample[int, int]{SampleSize: 20, ExpirePercent: 0.25}

		m.EXPECT().Size().Return(10)
		m.EXPECT().RandomKey().Return(0, false)

		sut.Expire(m)
	})

	t.Run("defaults sample size when zeroed", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		m := mockexpire.NewMockCacher[int, int](ctrl)
		sut := expire.RandomSample[int, int]{ExpirePercent: 0.25}

		m.EXPECT().Size().Return(0)

		sut.Expire(m)
		require.Equal(t, expire.DefaultSampleSize, sut.SampleSize)
	})

	t.Run("defaults expire percent when zeroed", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		m := mockexpire.NewMockCacher[int, int](ctrl)
		sut := expire.RandomSample[int, int]{SampleSize: 20}

		m.EXPECT().Size().Return(0)

		sut.Expire(m)
		require.Equal(t, expire.DefaultExpirePercent, sut.ExpirePercent)
	})
}
