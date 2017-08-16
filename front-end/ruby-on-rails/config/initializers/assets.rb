# Be sure to restart your server when you modify this file.

# Version of your assets, change this if you want to expire all your assets.
Rails.application.config.assets.version = '1.0'
Rails.application.config.assets.precompile += %w( drugDoseResponseCurve.js )
Rails.application.config.assets.precompile += %w( pieChart.js )
Rails.application.config.assets.precompile += %w( barPlot.js )
Rails.application.config.assets.precompile += %w( histogram.js )
Rails.application.config.assets.precompile += %w( venn.js )
Rails.application.config.assets.precompile += %w( waterfallPlot.js )
Rails.application.config.assets.precompile += %w( downloadSVG.js )
Rails.application.config.assets.precompile += %w( custom.js )
Rails.application.config.assets.precompile += %w( countup.js )
Rails.application.config.assets.precompile += %w( js-cookie.js )
Rails.application.config.assets.precompile += %w( js-cookie.js )
Rails.application.config.assets.precompile += %w( typeahead.js )
Rails.application.config.assets.precompile += %w( typed.js )
Rails.application.config.assets.precompile += %w( typeahead.js )
Rails.application.config.assets.precompile += %w( d3.v3.min.js )


# Add additional assets to the asset load path
# Rails.application.config.assets.paths << Emoji.images_path

# Precompile additional assets.
# application.js, application.css, and all non-JS/CSS in app/assets folder are already added.
# Rails.application.config.assets.precompile += %w( search.js )
