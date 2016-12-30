forecaster
==========

[![Build Status](https://travis-ci.org/gopher-finance/forecast.svg?branch=master)](https://travis-ci.org/gopher-finance/forecast)

Financial information retrieval and munging to automate security analysis on a large scale written in Go. Originally based in parts on [gofinance](https://github.com/aktau/gofinance).

Implementations of the Source interface can be found in the sources directory:

-	yahoofinance: implements **Source**. Queries **Yahoo Finance** for financial data, .
-	bloomberg: implements **Source**. Queries **Bloomberg** for financial data. (NOTE: Incomplete implementation)

Data Models:

-	models: data representation used by yahoofinance and bloomberg packages

The desktop cli app:

-	cmd/cli

Screenshot
==========

Everybody loves screenshots, even of console apps.

![Fancy screenshot of gofinance in action](http://aktau.github.io/gofinance/images/screen-01.png)

Requirements
============

-	Golang 1.7
-	[Glide](https://github.com/Masterminds/glide)

Install the app
---------------

If you haven't already done so, setup $GOPATH and $PATH:

```sh
# setup a $GOPATH
$ mkdir ~/go
$ export GOPATH=/my/home/path/go

# optional, add $GOPATH/bin to your $PATH so you don't have to navigate
# to it to run go binaries
$ export PATH=$PATH:$GOPATH/bin
```

After all that you can get the demonstration app (which will pull in all the necessary libraries automatically.

```sh
$ go get github.com/HT808s/gofinance/gofinance

# run the app, if you've added $GOPATH/bin to the PATH
$ make dep_install
$ make build
```

Run the app
===========

```sh
# check usage
$ ./bin/gofinance --help
# analyze Apple (default source=yahoo_yql)
$ ./bin/gofinance --symbols "AAPL"
# analyze Ferrari, Tesla and BMW
$ ./bin/gofinance -s "RACE" -s "TSLA" -s "BMW.DE"
# run with a different source
$ ./bin/gofinance -s "RACE" -s "TSLA" -s "BMW.DE" --source "yahoo_csv"
```

Features
========

-	Parallel fetching of data (with goroutines)
-	(Optional) caching of results (so the sources don't block you), configurable expiry time. This also allows for local pre-calculations that are too expensive to run on every fetch. (**Note** currently inhibited feature)

Todo
====

-	graph
	-	cli
		-	https://github.com/gizak/termui
-	look at (and possibly implement) some techniques from others packages:
	-	http://make.rafflecopter.com/side-project-jon-haskell.html
	-	[bateman](https://github.com/fearofcode/bateman)
	-	[tybee](https://github.com/yanatan16/tybee#design)
	-	[golang-spsa](https://github.com/yanatan16/golang-spsa)
-	classify the dividend level of stock, e.g.: aristocrat, ...
-	find a way to request echange rates (historical rates as well), to be able to indicate when's a good time to buy foreign securities. (it appears both Yahoo and Bloomberg have these as pseudo-securities:
	-	Yahoo: EURUSD=X -> http://chart.finance.yahoo.com/z?p=m50%2Cm200&q=l&s=VFEM.AS&t=2y&c=VYM,EURUSD=X)
	-	Bloomberg: EURUSD:CUR -> http://www.bloomberg.com/quote/EURUSD:CUR (already working)
	-	Alt: http://www.exchange-rates.org/history/EUR/USD/T
	-	StackExchange: http://quant.stackexchange.com/questions/141/what-data-sources-are-available-online
-	fallback to Yahoo Finance CSV when the YQL interface response with things like "too many instructions executed".
-	morningstar: scrape morningstar (all versions) as a ternary source, then we can be pretty sure that the app stays working for a while even if a source changes its format. Pull requests can then be submitted to fix the source. A couple of things Morningstar has that I haven't been able to find in others:
	-	It has industry and regional data for funds/ETF's, which could become very handy when factoring in currency risk and portfolio exposure.
	-	It lists the funds' own benchmark index and a Morningstar assigned one. We could make neat comparisons with those.
-	bloomberg: Funds/ETFs have a different layout than normal stocks on Bloomberg, adapt to that.
-	Combine data from Yahoo Finance and Bloomberg. (this should be implemented as a Source made out of multiple underlying Sources, like the Cache). If Yahoo Finance doesn't have the data on a company, Bloomberg usually has it. So there needs to be a function that determines if a Quote is "too incomplete". In the same vein, things like for example "dividend yield" can sometimes differ quite a bit between YF and Bloomberg, we need to find a way to reconcile this.
-	Combine data #2: sometimes neither Bloomberg nor Yahoo have the complete picture of a company. For example, Bloomberg usually has the full dividend yield information and Yahoo often seems to lack it. But Yahoo has more information in other departments, so this could be merged. To do this, we would need a way to specify for each datasource what it is intrinsically "good" at, such that when both Yahoo and Bloomberg give dividend information, Bloombergs' is preferred, et cetera.
-	Persist historical data locally (avoid getting blocked). This already happens for quotes if you query through a cache like the SqliteCache.
-	Add an optional ncurses-like userface, for example with termon: https://code.google.com/p/termon/
-	Plotting in the terminal, how to convert png -> ascii? Can't find easy libcaca bindings for now.
-	Somehow emulate the google finance stock screener: https://www.google.com/finance?ei=8EDhUuCpO4eHwAOklwE#stockscreener
-	Calculate dividend yield based on the price you actually paid (in aggregate) for stock you own. Why is this handy? For example: suppose you paid $20 for a stock that pays out $1 in dividends each year. It's dividend yield would be 5% (which is very good). If that stock rises to $40 the next year and the dividend payouts stay the same, the dividend yield would be 2.5% (which is still good but only half of the last time). However, this is not the yield you will be getting if you already own the stock (and bought it at $20). The yield would still be 5%, obviously. So, the yield as it stands is a good measure of what you're buying, but not of what you already bought.
-	Calculate effective yields and tax drag. This is hard to do (collect all the info), but once done it would be wonderful to determine how much a fund/stock/ETF is really going to net you. An example: the US levies a 30% tax on dividends as far as I know. This can be reduced to 15% if you or your broker make use of a double-tax treaty. Belgium has a 25% dividend tax which is (as of 2014) wholly non-negotiable, you always pay it, no matter where else you already paid taxes. In the Netherlands it is always 15%. In Ireland it appears to be 20% but I don't know if that counts for foreign companies who have their "domicile" in Ireland. It's all very confusing. At any rate, the most rosy tax climate I could personally get is first 15% (US) then 25% (BE). I'm too much of a realist to assume I will get that, but one has to take *a* value. So to calculate the effective yield I would use this formula: `yield * tax in source country * tax in receiving
	country`. In the case of for example a company with a dividend yield of 2.5%, that would become an effective yield of `2.5% * 0.85 * 0.75 =
	1.6%`. Which is suddenly a whole lot less great. As such, one realizes it takes quite a bit more dividend yield to beat a bank account by a nice margin. In a slightly worse case, I would also have to pay taxes in Holland, as that's where I usually buy shares (on the Amsterdam Euronext), so that would become `2.5% * 0.85 * 0.85 * 0.75 = 1.35%`. And that's how you hit rock bottom, even with an initially nice dividend yield. It gets even worse, if you want to see how your purchasing power will evolve, you have to take inflation into account. So let's do that. The average inflation in Belgium for 2013 was 1.11% (1.43% for Germany). That means that €1 in 2012 is worth `€1 / 1.0111
	= 0.989` in 2013. And so the final formula becomes `2.5% * 0.85 * 0.85
	-	0.75 - 1.11% = 0.24%`. And so, there's nothing left. Note that the levied taxes are massively important here. In the optimal case, where only Belgium and the US are paid (at the reduced US rate of 15%), we get an effective rate of 0.48%, which is not good but already double the other one. That said, it is often the case than when viewing yield estimates, source tax has already been incorporated! Let's take VUSA as an example, a Vanguard tracker for the S&P 500, domiciled in Ireland. The US-based variant of this is VOO. When you look at the yield difference between the two you'll notice that VOO has a yield of 2.09%, while VUSA has a yield of 1.63%. Why is this less? Because the dividend withholding tax has already been deducted (that's 15% because Ireland has a treaty with the US). Yet there's bound to be some extra cost, as a 15% tax rate would lead to a VUSA yield of 2.09%
	-	0.85 = 1.78%, but the reality says it's 1.63% (which corresponds to a "taxation" of 23%). This sounds bad for the european investor. But, there's a catch, the US yield is gross, so one has to deduct 30% WHT. Which puts the advantage in the European camp. But, every country wants another slice of the pie if you're in the EU. So you will in the end receive 23% * whatever_your_country_wants%. In my case, Belgium, that's 25%. So you lose again. Hopefully the EU will start acting up against these practices.
-	Dividend payout history (and other indicators) analysis. Even if the yield is good, it's possible that the company has just had a really bad run, and thus its share price is really low. This *might* mean that it's going to be less succesful (bad), or that the market is just in a slump (good). So companies that are just doing bad are to be avoiding, more so because they are unlikely to keep up their high dividend payments if they are on their way to the poor house. Some extra parameters gofinance could use to provide tips:
	-	historical dividends (are they growing, for how long?)
	-	strong growth and earnings
	-	no heavy drops in share price (we can disentangle the "losing company" case from the "bear market" case by comparing with the general direction of an index, be the index as representative as possible). Put shortly, fundamentals. One needs to make sure that the dividend is not going to be slashed and send the yield to kingdom come.
