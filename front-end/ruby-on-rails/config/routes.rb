Rails.application.routes.draw do

  get 'errors/not_found'

  get 'errors/internal_server_error'

	root 'static_pages#home'

	get 'search'       => 'search#index'

	# autocompletion
	get 'autocomplete' => 'autocomplete#suggest'

	# about
	get 'about' => 'profiles#about'

	# explore
	get 'explore' => 'profiles#explore'
	post 'explore' => 'profiles#explore'

	# cite_us
	get 'cite_us' => 'static_pages#cite_us'

	get 'pharmacogx' => 'static_pages#pharmacogx'

	# docs
	get 'docs'                       => 'docs#search'
	get 'docs/api'                   => 'docs#api'
	get 'docs/cell_line-vs-drug'     => 'docs#cell_drug'
	get 'docs/cell_lines'            => 'docs#cells'
	get 'docs/tissues'               => 'docs#tissues'
	get 'docs/drugs'                 => 'docs#drugs'
	get 'docs/datasets'              => 'docs#datasets'
	get 'docs/datasets-intersection' => 'docs#datasets_intersection'
	get 'docs/search'                => 'docs#search'
	get 'docs/videos'                => 'docs#videos'

	# profile pages
	get 'cell_lines'   => 'profiles#cell_lines'
	get 'tissues'      => 'profiles#tissues'
	get 'drugs'        => 'profiles#drugs'
	get 'datasets'     => 'profiles#datasets'
	get 'experiments'  => 'profiles#experiments'
	get 'targets'  => 'profiles#targets'
	get 'contact_us'   => 'profiles#contact_us'
	post 'contact_us'   => 'profiles#contact_us'

	get 'cell_lines/:id'       => 'profiles#cell_lines'
	get 'tissues/:id'          => 'profiles#tissues'
	get 'drugs/:id'            => 'profiles#drugs'
	get 'datasets/:id'         => 'profiles#datasets'
	get 'targets/:id'         => 'profiles#targets'

  get 'download' => 'static_pages#download'

  # errors
  match "/404", :to => "errors#not_found", :via => :all
  match "/500", :to => "errors#internal_server_error", :via => :all

	# For details on the DSL available within this file, see http://guides.rubyonrails.org/routing.html
end
